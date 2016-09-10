package dbsrv

import (
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/reform"
)

const (
	DaySeconds int64 = 3600 * 24
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
	for _, iorder := range orders {
		order := iorder.(*front.Order)
		cols, err := dbs.OrderMaintanence(order)
		if err != nil {
			return err
		}
		if len(cols) != 0 {
			if err = dbs.GetDB().UpdateColumns(order, append(cols, "State")...); err != nil {
				return err
			}
		}
	}
	return nil
}

// exclude State
func (dbs *DbService) OrderMaintanence(order *front.Order) (cols []string, err error) {
	if order.PayPoints != 0 {
		return
	}

	db := dbs.GetDB()
	tnow := time.Now()
	now := tnow.Unix()

	if dbs.IsOrderAutoCompletedUnsaved(order) {
		order.CompletedAt = order.DeliveredAt + int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds
		if order.EvalStartedAt != 0 && order.EvalStartedAt < order.CompletedAt {
			order.CompletedAt = order.EvalStartedAt
		} else if order.EvalAt != 0 && order.EvalAt < order.CompletedAt {
			order.CompletedAt = order.EvalAt
		} else {
			order.AutoCompleted = true
		}
		order.State = front.TOrderStateCompleted
		cols = append(cols, "AutoCompleted", "CompletedAt")
	}
	if dbs.IsOrderAutoEvaledUnsaved(order) {
		if order.CompletedAt == 0 {
			order.CompletedAt = order.DeliveredAt + int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds
			if order.CompletedAt > order.EvalAt {
				order.CompletedAt = order.EvalAt
			}
			cols = append(cols, "CompletedAt")
		}
		order.AutoEvaled = true
		order.EvalAt = order.CompletedAt + int64(dbs.config.Order.EvalTimeoutDay)*DaySeconds
		order.State = front.TOrderStateEvaled
		cols = append(cols, "AutoEvaled", "EvalAt")
	}
	if dbs.IsOrderHistoryUnsaved(order) {
		order.HistoryAt = order.EvalAt + int64(dbs.config.Order.HistoryTimeoutDay)*DaySeconds
		order.State = front.TOrderStateHistory
		cols = append(cols, "HistoryAt")
		if order.EvalAt == 0 {
		}
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
					toStoreMap[item.StoreID] = item.Price * item.Quantity
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
			ds := dbs.DS.Where(goqu.I(front.StoreTable.PK()).Eq(id)).Limit(1)
			var count uint64
			count, err = db.DsCount(front.StoreTable, ds)
			if err != nil {
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
			ds := dbs.DS.Where(goqu.I(models.UserTable.PK()).Eq(id)).Limit(1)
			var count uint64
			count, err = db.DsCount(models.UserTable, ds)
			if err != nil {
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
		var usr1CashBalance int
		total := order.PayAmount - order.DeliverFee
		toUsr1 := total * dbs.config.Money.User1RebatePercent / 100
		if order.User1 != 0 {
			ds := dbs.DS.Where(goqu.I(models.UserTable.PK()).Eq(order.User1))
			var count uint64
			count, err = db.DsCount(models.UserTable, ds)
			if err != nil {
				return
			}
			if count == 1 {
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
			CreatedAt:  now,
			Type:       front.TUserCashTrade,
			OrderTotal: order.PayAmount,
			OrderID:    order.ID,
			Amount:     toPlatform,
			Balance:    top.Balance + toPlatform,
		})
		if err != nil {
			return
		}

		// ABC
		var abcItem *front.OrderItem
		if len(items) == 1 && items[0].IsABC {
			abcItem = items[0]

			// 1. get user
			var usr models.User
			if err = db.FindByPrimaryKeyTo(&usr, order.UserID); err != nil {
				return
			}

			// 2. get exist vips of user
			var iUserVips []reform.Struct
			ds = dbs.DS.Where(goqu.I("$UserID").Eq(order.UserID), goqu.I("$ExpiresAt").Gt(now))
			iUserVips, err = db.DsSelectAllFrom(front.VipRebateOriginTable, ds)
			if err != nil {
				return
			}

			// 3. find valid/last vip
			var userVips []*front.VipRebateOrigin
			var userCurrentVip *front.VipRebateOrigin
			var userLastVip *front.VipRebateOrigin
			for _, ivip := range iUserVips {
				vip := ivip.(*front.VipRebateOrigin)
				userVips = append(userVips, vip)
				if vip.Valid(now) {
					userCurrentVip = vip
				}
				if userLastVip == nil || userLastVip.NotBefore < vip.NotBefore {
					userLastVip = vip
				}
			}

			// 4. save VipRebateOrigin of user
			userVip := front.VipRebateOrigin{
				UserID:    order.UserID,
				CreatedAt: order.CreatedAt,
				OrderID:   order.ID,
				ItemID:    abcItem.ID,
				Amount:    abcItem.Price,
				Balance:   abcItem.Price,
				User1:     order.User1,
				//				userVip.NotBefore = now
				//				userVip.ExpiresAt = nextYearNow
				//				userVip.User1Used = false
			}

			if userCurrentVip == nil {
				userVip.NotBefore = now
				userVip.ExpiresAt = tnow.AddDate(1, 0, 0).Unix()
			} else {
				latest := time.Unix(userLastVip.NotBefore, 0)
				userVip.NotBefore = latest.AddDate(1, 0, 0).Unix()
				userVip.ExpiresAt = latest.AddDate(2, 0, 0).Unix()
			}
			if err = db.Insert(&userVip); err != nil {
				return
			}

			// 5. thaw some cash if exsit
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
					var top front.UserCash
					ds = dbs.DS.Where(goqu.I("$UserID").Eq(order.UserID)).Order(goqu.I("$CreatedAt").Desc())
					if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
						return
					}

					// to rebate cash? impossible!
					err = db.Insert(&front.UserCash{
						UserID:    order.UserID,
						CreatedAt: now,
						Type:      item.Type,
						Amount:    int(item.Amount),
						Balance:   top.Balance + int(item.Amount),
						OrderID:   item.OrderID,
						Remark:    item.Remark,
					})
					if err != nil {
						return
					}
				}
			}

			// 6. move rebate of user1 to frontend.
			// backend just accepts the choice of user1
		}
	}

	return cols, nil
}

// OrderCompleted
func (dbs *DbService) IsOrderCompleted(order *front.Order) bool {
	return order.DeliveredAt != 0 && (order.CompletedAt != 0 ||
		order.HistoryAt != 0 || order.EvalStartedAt != 0 || order.EvalAt != 0 ||
		time.Now().Unix()-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds)
}

func (dbs *DbService) IsOrderAutoCompletedUnsaved(order *front.Order) bool {
	return order.DeliveredAt != 0 && order.CompletedAt == 0 &&
		time.Now().Unix()-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds
}

// OrderEvaled
func (dbs *DbService) IsOrderEvaled(order *front.Order) bool {
	return order.DeliveredAt != 0 &&
		(order.EvalAt != 0 || order.HistoryAt != 0 || dbs.isOrderAutoEvaledUnsaved(order))
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
