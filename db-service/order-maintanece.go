package dbsrv

import (
	"time"

	"gopkg.in/doug-martin/goqu.v3"
	"gopkg.in/reform.v1"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

func (dbs *DbService) OrdersMaintanence() error {
	now := time.Now().Unix()
	ds := dbs.DS.Where(
		goqu.I("$DeliveredAt").Gt(0),
		goqu.I("$DeliveredAt").Lte(now-int64(dbs.config.Order.CompleteTimeoutDay)*3600*24),
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
		order.CompletedAt = order.DeliveredAt + int64(dbs.config.Order.CompleteTimeoutDay)*3600*24
		order.State = front.TOrderStateCompleted
		cols = append(cols, "AutoCompleted", "CompletedAt", "State")
	}
	if dbs.IsOrderAutoEvaledUnsaved(order) {
		order.AutoEvaled = true
		order.EvalAt = order.CompletedAt + int64(dbs.config.Order.EvalTimeoutDay)*3600*24
		order.State = front.TOrderStateEvaled
		cols = append(cols, "AutoEvaled", "EvalAt", "State")
	}
	if dbs.IsOrderHistoryUnsaved(order) {
		order.HistoryAt = order.EvalAt + int64(dbs.config.Order.HistoryTimeoutDay)*3600*24
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

	if !order.Rebated {
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
		total := order.PayAmount - order.DeliverFee
		toUsr1 := total * dbs.config.Money.User1RebatePercent / 100
		if order.User1 != 0 {
			ds := dbs.DS.Where(goqu.I(models.UserTable.PK()).Eq(order.User1))
			count, err1 := db.DsCount(models.UserTable, ds)
			err = db.FindByPrimaryKeyTo(&usr1, order.User1)
			if err == nil {
				var top front.UserCash
				ds := dbs.DS.Where(goqu.I("$UserID").Eq(order.User1)).Order(goqu.I("$CreatedAt").Desc())
				if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
					return
				}

				err = db.Insert(&front.UserCash{
					UserID:    order.User1,
					CreatedAt: now,
					Type:      front.TUserCashRebate,
					Amount:    int(toUsr1),
					Balance:   top.Balance + int(toUsr1),
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

		// ABCs
		if order.User1 != 0 {
			var abcs []*front.OrderItem
			for _, item := range items {
				if item.IsABC {
					abcs = append(abcs, item)
				}
			}
			if len(abcs) > 0 {
				// TODO multi abcs in user1
				// 1. VipRebateOrigin of user1 from the last year
				// 2. set vip to user if needed
				// 3. start rebate of user1 if counter is enough
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
		time.Now().Unix()-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay)*3600*24)
}

func (dbs *DbService) IsOrderAutoCompletedUnsaved(order *front.Order) bool {
	return order.DeliveredAt != 0 && order.CompletedAt == 0 &&
		time.Now().Unix()-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay)*3600*24
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
		return now-order.CompletedAt > int64(dbs.config.Order.EvalTimeoutDay)*3600*24
	}
	return now-order.DeliveredAt > int64(dbs.config.Order.CompleteTimeoutDay+dbs.config.Order.EvalTimeoutDay)*3600*24
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
		return now-order.EvalAt > int64(dbs.config.Order.HistoryTimeoutDay)*3600*24
	}

	if order.CompletedAt != 0 {
		return now-order.CompletedAt > int64(dbs.config.Order.HistoryTimeoutDay+dbs.config.Order.EvalTimeoutDay)*3600*24
	}

	return now-order.DeliveredAt > int64(dbs.config.Order.HistoryTimeoutDay+
		dbs.config.Order.EvalTimeoutDay+dbs.config.Order.CompleteTimeoutDay)*3600*24
}
