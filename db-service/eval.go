package dbsrv

import (
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/lok"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/reform"
)

func (dbs *DbService) EvalSave(tokUsr *models.User, orderId, itemId uint, payload *front.EvalItem) (*front.EvalResponse, error) {
	name := []rune(tokUsr.Nickname)
	switch l := len(name); l {
	case 1:
	case 2:
		name[1] = '*'
	case 3:
		name[1] = '*'
		name[2] = '*'
	default:
		for i := 1; i < l-1; i++ {
			name[i] = '*'
		}
	}

	payload.EvalName = string(name)
	payload.EvalAt = time.Now().Unix()

	if !lok.OrderLok.Lock(orderId) {
		return nil, cerr.OrderTmpLocked
	}
	defer lok.OrderLok.Unlock(orderId)

	var order front.Order
	var ra uint
	err := dbs.InTx(func(tx *DbService) error {
		db := tx.GetDB()

		ds := dbs.DS.Where(goqu.I(front.OrderTable.PK()).Eq(orderId)).Where(goqu.I("$UserID").Eq(tokUsr.ID))
		err := db.DsSelectOneTo(&order, ds)
		if err != nil {
			return err
		}

		// TOrderStateEvalStarted included
		if err := PermitOrderState(&order, front.TOrderStateEvaled); err != nil {
			return cerr.NoWayToTargetState
		}

		if dbs.IsEvalTimeout(&order) {
			return cerr.OrderEvalTimeout
		}

		// because payload is EvalItem
		ds = dbs.DS.Where(goqu.I(front.OrderItemTable.ToCol("UserID")).Eq(tokUsr.ID))
		if itemId == 0 {
			ds = ds.Where(goqu.I(front.OrderItemTable.ToCol("OrderID")).Eq(orderId)).Where(goqu.I("$EvalAt").IsNull())
		} else {
			ds = ds.Where(goqu.I(front.OrderItemTable.PK()).Eq(itemId))
		}

		ra, err = db.DsUpdateStruct(payload, ds)
		if err != nil {
			return err
		}

		var evalOrder bool
		if itemId == 0 || order.EvalAt != 0 {
			evalOrder = true
		} else {
			ds = ds.Where(goqu.I("$OrderID").Eq(orderId)).Where(goqu.I("$EvalAt").IsNull())
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
				cols = append(cols, "State", "EvalAt")
			}
			if evalStarted {
				order.EvalStartedAt = payload.EvalAt
				cols = append(cols, "EvalStartedAt")
				if order.State != front.TOrderStateEvaled {
					order.State = front.TOrderStateEvalStarted
					cols = append(cols, "State")
				}
			}
			return db.UpdateColumns(&order, cols...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &front.EvalResponse{
		Order:    &order,
		Evaled:   ra,
		EvalAt:   payload.EvalAt,
		EvalName: payload.EvalName,
	}, nil
}
