package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type orderChangeStateInput struct {
	order   *front.Order
	tokUsr  *models.User
	payload *front.OrderChangeStatePayload
	chanErr <-chan error
}

func (hub *OrderHub) OrderChangeState(
	order *front.Order, tokUsr *models.User, payload *front.OrderChangeStatePayload,
) (err error) {
	chanErr := make(chan error)
	in := &orderChangeStateInput{order, tokUsr, payload, chanErr}
	hub.chanInput <- in

	err = <-chanErr
	return
}

func (hub *OrderHub) onOrderChangeState(in *orderChangeStateInput) {
	in.chanErr <- hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.PayOrder(in.order, in.tokUsr, in.payload)
		return
	})
}
