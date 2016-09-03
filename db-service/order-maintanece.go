package dbsrv

import (
	"time"

	"gopkg.in/doug-martin/goqu.v3"
	"gopkg.in/reform.v1"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

const (
	DaySeconds uint64 = 3600 * 24
)

func (dbs *DbService) OrdersMaintanence() error {
	now := time.Now().Unix()
	ds := dbs.DS.Where(
		goqu.I("$DeliveredAt").Gt(0),
		goqu.I("$DeliveredAt").Lte(now-int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds),
		goqu.Or(
			goqu.I("$CompletedAt").Eq(0),
			goqu.I("$EvalAt").Eq(0),
			goqu.I("$EvalStartedAt").Eq(0),
			goqu.I("$HistoryAt").Eq(0),
		),
	)
	orders, err := dbs.GetDB().DsFindAllFrom(front.OrderTable, ds)
	if err != nil {
		return err
	}
	for _, order := range orders {
		if err = dbs.OrderMaintanence(order.(*front.Order)); err != nil {
			return err
		}
	}
	return nil
}

func (dbs *DbService) OrderMaintanence(order front.Order) (changed *front.Order, err error) {
	db := dbs.GetDB()
	now := time.Now().Unix()

	var cols []string

	if dbs.IsOrderAutoCompletedUnsaved(order) {
		order.AutoCompleted = true
		order.CompletedAt = order.DeliveredAt + int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds
		order.State = front.TOrderStateCompleted
		cols = append(cols, "AutoCompleted", "CompletedAt", "State")
	}
	if dbs.IsOrderAutoEvaledUnsaved(order) {
		order.AutoEvaled = true
		order.EvalAt = order.CompletedAt + int64(dbs.config.Order.EvalTimeoutDay)*DaySeconds
		order.State = front.TOrderStateEvaled
		cols = append(cols, "AutoEvaled", "EvalAt", "State")
	}
	if dbs.IsOrderHistoryUnsaved(order) {
		order.HistoryAt = order.EvalAt + int64(dbs.config.Order.HistoryTimeoutDay)*DaySeconds
		order.State = front.TOrderStateHistory
		cols = append(cols, "HistoryAt", "State")
	}

	items, err1 := dbs.GetOrderItems(order)
	if err1 != nil {
		err = err1
		return
	}
	if items == nil {
		err = cerr.OrderItemNotFound
		return
	}

	if order.CompletedAt != 0 && !order.Rebated {
		cols = append(cols, "Rebated")
		order.Rebated = true
		// 1. split money to parts
		toStoreMap := make(map[uint]uint)
		toStore1Map := make(map[uint]uint)
		for _, item := range items {
			if item.StoreID != 0 {
				if to, ok := toStoreMap[item.StoreID]; ok {
					toStoreMap[item.StoreID] = to + item.Price*item.Quantity
				} else {
					toStoreMap[item.Store1] = item.Price * item.Quantity
				}
			}
			if item.Store1 != 0 {
				if to, ok := toStore1Map[item.Store1]; ok {
					toStore1Map[item.Store1] = to + item.Price*item.Quantity
				} else {
					toStore1Map[item.Store1] = item.Price * item.Quantity
				}
			}
		}

		// 2. store
		var toStoreAll uint
		for id, to := range toStoreMap {
			ds := dbs.DS.Where(goqu.I(front.StoreTable.PK()).Eq(id))
			count, err1 := db.DsCount(front.StoreTable, ds)
			if err1 != nil {
				err = err1
				return
			}
			if count != 0 {
				toStore := to * (100 - dbs.config.Money.StoreSaleFeePercent) / 100
				toStoreAll += toStore

				var top front.StoreCash
				ds := dbs.DS.Where(goqu.I("$StoreID").Eq(id)).Order(goqu.I("$CreatedAt").Desc())
				if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
					return
				}

				err = db.Insert(&front.StoreCash{
					StoreID:   id,
					CreatedAt: now,
					OrderID:   order.ID,
					Amount:    toStore,
					Balance:   top.Balance + toStore,
				})
				if err != nil {
					return
				}
			}
		}

		// 3. store1
		var toStore1All uint
		for id, to := range toStoreMap {
			ds := dbs.DS.Where(goqu.I(models.UserTable.PK()).Eq(id))
			count, err1 := db.DsCount(models.UserTable, ds)
			if err1 != nil {
				err = err1
				return
			}
			if count != 0 {
				toStore1 := to * dbs.config.Money.Store1RebatePercent / 100
				toStore1All += toStore1

				var top front.UserCash
				ds := dbs.DS.Where(goqu.I("$StoreID").Eq(id)).Order(goqu.I("$CreatedAt").Desc())
				if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
					return
				}

				err = db.Insert(&front.UserCash{
					UserID:    id,
					CreatedAt: now,
					Type:      front.TUserCashStoreRebate,
					Amount:    int(toStore1),
					Balance:   top.Balance + int(toStore1),
					OrderID:   order.ID,
				})
				if err != nil {
					return
				}
			}
		}

		// 4. user1
		var usr1 models.User
		var usr1CashBalance int
		total := order.PayAmount - order.DeliverFee
		toUsr1 := total * dbs.config.Money.User1RebatePercent / 100
		if order.User1 != 0 {
			//			ds := dbs.DS.Where(goqu.I(models.UserTable.PK()).Eq(order.User1))
			//			count, err1 := db.DsCount(models.UserTable, ds)
			err = db.FindByPrimaryKeyTo(&usr1, order.User1)
			if err == nil {
				var top front.UserCash
				ds := dbs.DS.Where(goqu.I("$UserID").Eq(order.User1)).Order(goqu.I("$CreatedAt").Desc())
				if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
					return
				}
				usr1CashBalance = top.Balance + int(toUsr1)

				err = db.Insert(&front.UserCash{
					UserID:    order.User1,
					CreatedAt: now,
					Type:      front.TUserCashRebate,
					Amount:    int(toUsr1),
					Balance:   usr1CashBalance,
					OrderID:   order.ID,
				})
				if err != nil {
					return
				}
			} else if err != reform.ErrNoRows {
				return
			}
		}

		// 5. platform
		toPlatform := order.PayAmount - toStoreAll - toStore1All - toUsr1
		var top models.PlatformCash
		ds := dbs.DS.Order(goqu.I("$CreatedAt").Desc())
		if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
			return
		}

		err = db.Insert(&models.PlatformCash{
			CreatedAt: now,
			Type:      front.TUserCashTrade,
			Total:     order.PayAmount,
			OrderID:   order.ID,
			Amount:    toPlatform,
			Balance:   top.Balance + toPlatform,
		})
		if err != nil {
			return
		}

		// ABC
		var abcItem *front.OrderItem
		if len(items) == 1 && items[0].IsABC {
			abcItem = items[0]
			// 1. save VipRebateOrigin to user
			err = db.Insert(&models.VipRebateOrigin{
				UserID:    order.UserID,
				CreatedAt: order.CreatedAt,
				OrderID:   order.ID,
				ItemID:    abcItem.ID,
				Amount:    abcItem.Price,
				Balance:   abcItem.Price,
				User1:     order.User1,
			})
			if err != nil {
				return
			}

			var usr models.User
			if err = db.FindByPrimaryKeyTo(&usr, order.UserID); err != nil {
				return
			}

			// 2. set vip of user
			tnow := time.Now()
			year, month, day := tnow.Date()
			begin := time.Date(year, month, day, tnow.Hour(), tnow.Minute(), tnow.Second(), 0, time.Local).Unix()
			if usr.VipAt < begin {
				// expires, use NextVipAt as VipAt
				usr.VipAt, usr.NextVipAt = usr.NextVipAt, 0
			}
			if usr.VipAt < begin {
				// expires, set vip
				usr.VipAt = order.CompletedAt
			} else {
				// not expires, set next vip
				usr.NextVipAt = order.CompletedAt
			}
			if err = db.UpdateColumns(&usr, "VipAt", "NextVipAt"); err != nil {
				return
			}

			// 3. thaw some cash if exsit
			ds = dbs.DS.Where(goqu.I("$UserID").Eq(order.UserID), goqu.I("$ThawedAt").Eq(0))
			var freeze []reform.Struct
			freeze, err = db.DsSelectAllFrom(front.UserCashFrozenTable, ds)
			if err != nil {
				return
			}
			if len(freeze) > 0 {
				_, err = db.DsUpdateColumns(&front.UserCashFrozen{ThawedAt: now}, ds, "ThawedAt")
				if err != nil {
					return
				}

				for _, itemi := range freeze {
					item := itemi.(*front.UserCashFrozen)
					if item.Stages == 0 {
						var top front.UserCash
						ds = dbs.DS.Where(goqu.I("$UserID").Eq(order.UserID)).Order(goqu.I("$CreatedAt").Desc())
						if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
							return
						}

						err = db.Insert(&front.UserCash{
							UserID:    order.UserID,
							CreatedAt: now,
							Type:      item.Type,
							Amount:    int(item.Amount),
							Balance:   top.Balance + int(item.Amount),
							OrderID:   item.OrderID,
							Remark:    item.Remark,
						})
					} else {
						err = db.Insert(&front.UserCashRebate{
							UserID:    order.UserID,
							OrderID:   item.OrderID,
							CreatedAt: now,
							Type:      item.Type,
							Amount:    item.Amount,
							Remark:    item.Remark,
							Stages:    item.Stages,
							DoneAt:    0,
						})
					}
					if err != nil {
						return
					}
				}
			}

			// 4. rebate for user1
			if order.User1 != 0 {
				usr1WasVip := usr1.VipAt != 0
				if usr1.VipAt < begin {
					usr1.VipAt, usr1.NextVipAt = usr1.NextVipAt, 0
				}
				if usr1.VipAt < begin {
					// not vip
					usr1.VipAt = 0
					if usr1WasVip {
						// 4.1 rebate counter to user1 cash if user1 is not vip but was right now
						// we DO NOT record order_id
						if usr1.NotRebatedCounter > 1 {
							log.Errorln("NotRebatedCounter err:", usr1.NotRebatedCounter)
						}
						for i := 0; i < usr1.NotRebatedCounter; i++ {
							usr1CashBalance += int(dbs.config.Money.RewardFromVipCent)
							err = db.Insert(&front.UserCash{
								UserID:    order.User1,
								CreatedAt: now,
								Type:      front.TUserCashReward,
								Amount:    int(dbs.config.Money.RewardFromVipCent),
								Balance:   usr1CashBalance,
								Remark:    "VIP expires auto reward",
							})
							if err != nil {
								return
							}
						}
						// set 1 to new vip lifecycle
						usr1.NotRebatedCounter = 1
					} else {
						// 4.2 just add to NotRebatedCounter if user1 is not vip
						usr1.NotRebatedCounter++
					}
					// not vip end
				} else {
					// 4.3 user1 is vip, so we can check if needing rebate or reward
					// get current VipRebateOrigin
					ds = dbs.DS.Where(goqu.I("$UserID").Eq(order.User1), goqu.I("$CreatedAt").Eq(usr1.VipAt))
					var rebateOrigin models.VipRebateOrigin
					err = db.DsSelectOneTo(&rebateOrigin, ds)
					if err != nil {
						return
					}

					if rebateOrigin.Balance == 0 {
						// 4.4 reward to user1
						usr1CashBalance += int(dbs.config.Money.RewardFromVipCent)
						err = db.Insert(&front.UserCash{
							UserID:    order.User1,
							CreatedAt: now,
							Type:      front.TUserCashReward,
							Amount:    int(dbs.config.Money.RewardFromVipCent),
							Balance:   usr1CashBalance,
							OrderID:   order.ID,
						})
						if err != nil {
							return
						}
					} else {
						// 4.5 rebate to user1
						usr1.NotRebatedCounter++
						if usr1.NotRebatedCounter >= 2 {
							usr1.NotRebatedCounter -= 2
							// 4.5.1 find all VipRebateOrigin from next level users of user1
							ds := dbs.DS.Where(goqu.I("$User1").Eq(order.User1), goqu.I("$User1Used").IsNotTrue())
							if rebateOrigin.Balance == rebateOrigin.Amount {
								toRebate := rebateOrigin.Amount / 2
								rebateOrigin.Balance = rebateOrigin.Amount - toRebate
							}
						}
					}

				}

			}
		}

		if order.User1 != 0 {
			if err = db.Save(&usr1); err != nil {
				return
			}
		}
	}

	if len(cols) != 0 {
		if err = db.UpdateColumns(&order, cols...); err == nil {
			changed = &order
		}
	}
	return
}

// OrderCompleted
func (dbs *DbService) IsOrderCompleted(order *front.Order) bool {
	return order.DeliveredAt != 0 && (order.HistoryAt != 0 ||
		order.CompletedAt != 0 || order.EvalStartedAt != 0 || order.EvalAt != 0 ||
		time.Now().Unix()-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds)
}

func (dbs *DbService) IsOrderAutoCompletedUnsaved(order *front.Order) bool {
	return order.DeliveredAt != 0 && order.CompletedAt == 0 &&
		time.Now().Unix()-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds
}

// OrderEvaled
func (dbs *DbService) IsOrderEvaled(order *front.Order) bool {
	return order.DeliveredAt != 0 &&
		(order.HistoryAt != 0 || order.EvalAt != 0 || dbs.isOrderAutoEvaledUnsaved(order))
}

func (dbs *DbService) IsOrderAutoEvaledUnsaved(order *front.Order) bool {
	return order.DeliveredAt != 0 && dbs.isOrderAutoEvaledUnsaved(order)
}

func (dbs *DbService) isOrderAutoEvaledUnsaved(order *front.Order) bool {
	now := time.Now().Unix()
	if order.CompletedAt != 0 {
		return now-order.CompletedAt > int64(dbs.config.Order.EvalTimeoutDay)*DaySeconds
	}
	return now-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay+dbs.config.Order.EvalTimeoutDay)*DaySeconds
}

// OrderHistory
func (dbs *DbService) IsOrderHistory(order *front.Order) bool {
	return order.DeliveredAt != 0 && (order.HistoryAt != 0 || dbs.isOrderHistoryUnsaved(order))
}

func (dbs *DbService) IsOrderHistoryUnsaved(order *front.Order) bool {
	return order.DeliveredAt != 0 && order.HistoryAt == 0 && dbs.isOrderHistoryUnsaved(order)
}

func (dbs *DbService) isOrderHistoryUnsaved(order *front.Order) bool {
	now := time.Now().Unix()
	if order.EvalAt != 0 {
		return now-order.EvalAt > int64(dbs.config.Order.HistoryTimeoutDay)*DaySeconds
	}

	if order.CompletedAt != 0 {
		return now-order.CompletedAt > int64(dbs.config.Order.HistoryTimeoutDay+dbs.config.Order.EvalTimeoutDay)*DaySeconds
	}

	return now-order.DeliveredAt > int64(dbs.config.Order.HistoryTimeoutDay+
		dbs.config.Order.EvalTimeoutDay+dbs.config.Order.CompleteTimeoutDay)*DaySeconds
}
