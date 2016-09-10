package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type orderEvalInput struct {
	order   *front.Order
	ra      *uint
	tokUsr  *models.User
	orderId uint
	itemId  uint
	payload *front.EvalItem
	chanErr <-chan error
}

func (hub *OrderHub) EvalSave(
	order *front.Order, ra *uint, tokUsr *models.User, orderId, itemId uint, payload *front.EvalItem,
) (err error) {
	chanErr := make(chan error)
	in := &orderEvalInput{order, ra, tokUsr, orderId, itemId, payload, chanErr}
	hub.chanInput <- in

	err = <-chanErr
	return
}

func (hub *OrderHub) onEvalSave(in *orderEvalInput) {
	in.chanErr <- hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.EvalSave(order, ra, tokUsr, orderId, itemId, payload)
		return
	})
}
