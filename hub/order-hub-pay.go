package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type payOrderInput struct {
	tokUsr    *models.User
	payload   *front.OrderPayPayload
	chanOrder <-chan *front.Order
	chanErr   <-chan error
}

func (hub *OrderHub) PayOrder(tokUsr *models.User, payload *front.OrderPayPayload) (o *front.Order, err error) {
	chanOrder := make(chan *front.Order)
	chanErr := make(chan error)
	in := &payOrderInput{tokUsr, payload, chanOrder, chanErr}

	hub.chanInput <- in

	select {
	case o = <-chanOrder:
	case err = <-chanErr:
	}
	return
}

func (hub *OrderHub) onPayOrder(in *payOrderInput) {
	var order *front.Order
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		order, err = tx.PayOrder(in.tokUsr, in.payload)
	})
	if err != nil {
		in.chanErr <- err
	} else {
		in.chanResult <- order
	}
}
