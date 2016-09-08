package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
)

type prepayOrderResult struct {
	order   *front.Order
	prepaid bool
}

type prepayOrderInput struct {
	userId     uint
	orderId    uint
	chanResult <-chan *prepayOrderResult
	chanErr    <-chan error
}

func (hub *OrderHub) PrepayOrder(userId, orderId uint) (order *front.Order, prepaid bool, err error) {
	chanResult := make(chan *prepayOrderResult)
	chanErr := make(chan error)
	in := prepayOrderInput{userId, orderId, chanResult, chanErr}
	hub.chanInput <- in

	select {
	case result := <-chanResult:
		order, prepaid = result.order, result.prepaid
	case err = <-chanErr:
	}
	return
}

func (hub *OrderHub) onPrepayOrder(tx *dbsrv.DbService, in *prepayOrderInput) {
	var order *front.Order
	var prepaid bool
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		order, prepaid, err = tx.PrepayOrder(in.PrepayOrderInput)
	})
	if err != nil {
		in.chanErr <- err
	} else {
		in.chanResult <- &prepayOrderResult{order, prepaid}
	}
}
