package dbsrv

import (
	"gopkg.in/doug-martin/goqu.v3"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/reform"
)

func (dbs *DbService) EvalSave(
	order *front.Order, ra *uint, tokUsr *models.User, orderId, itemId uint, payload *front.EvalItem,
) error {

	db := dbs.GetDB()

	ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(orderId)).Where(goqu.I("$UserID").Eq(tokUsr.ID))
	err := db.DsSelectOneTo(order, ds)
	if err != nil {
		return err
	}

	// TOrderStateEvalStarted included
	if err := PermitOrderState(order, front.TOrderStateEvaled); err != nil {
		return cerr.NoWayToTargetState
	}

	if dbs.IsOrderHistory(order) {
		return cerr.OrderEvalTimeout
	}

	// because payload is EvalItem
	ds = dbs.DS.Where(goqu.I(front.OrderItemTable.ToCol("UserID")).Eq(tokUsr.ID))
	if itemId == 0 {
		ds = ds.Where(goqu.I(front.OrderItemTable.ToCol("OrderID")).Eq(orderId)).Where(goqu.I("$EvalAt").Eq(0))
	} else {
		ds = ds.Where(goqu.I(front.OrderItemTable.PK()).Eq(itemId))
	}

	*ra, err = db.DsUpdateStruct(payload, ds)
	if err != nil {
		return err
	}

	var evalOrder bool
	if itemId == 0 || order.EvalAt != 0 {
		evalOrder = true
	} else {
		ds = ds.Where(goqu.I("$OrderID").Eq(orderId)).Where(goqu.I("$EvalAt").Eq(0))
		unevaled, err := db.DsCount(front.OrderItemTable, ds)
		if err != nil && err != reform.ErrNoRows {
			return err
		}
		evalOrder = unevaled == 0
	}

	evalStarted := order.EvalStartedAt == 0
	if evalOrder || evalStarted {
		var cols []string
		if evalOrder {
			order.State = front.TOrderStateEvaled
			order.EvalAt = payload.EvalAt
			cols = append(cols, "EvalAt")
		}
		if evalStarted {
			order.EvalStartedAt = payload.EvalAt
			cols = append(cols, "EvalStartedAt")
			if order.State != front.TOrderStateEvaled {
				order.State = front.TOrderStateEvalStarted
			}
		}
		if order.CompletedAt == 0 {
			order.CompletedAt = order.DeliveredAt + int64(dbs.config.Order.CompleteTimeoutDay)*DaySeconds
			if order.CompletedAt > payload.EvalAt {
				order.CompletedAt = payload.EvalAt
			}
			cols = append(cols, "CompletedAt")

			// TODO prove it
			cols2, err := dbs.OrderMaintanence(order)
			if err != nil {
				return err
			}
			cols = append(cols, cols2...)
		}
		return db.UpdateColumns(order, append(cols, "State")...)
	}
	return nil
}
