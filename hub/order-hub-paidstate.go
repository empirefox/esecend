package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type orderPaidStateInput struct {
	tokUsr    *models.User
	orderId   uint
	chanOrder <-chan *front.Order
	chanErr   <-chan error
}

func (hub *OrderHub) OrderPaidState(tokUsr *models.User, orderId uint) (order *front.Order, err error) {
	chanOrder := make(chan *front.Order)
	chanErr := make(chan error)
	in := &orderPaidStateInput{tokUsr, orderId, chanOrder, chanErr}
	hub.chanInput <- in

	select {
	case order = <-chanOrder:
	case err = <-chanErr:
	}
	return
}

func (hub *OrderHub) onOrderPaidState(in *orderPaidStateInput) {
	var order *front.Order
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		order, err = tx.OrderPaidState(in.tokUsr, in.orderId)
		return
	})
	if err != nil {
		in.chanErr <- err
	} else {
		in.chanResult <- order
	}
}
