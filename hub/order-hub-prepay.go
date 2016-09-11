package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type prepayOrderResult struct {
	order *front.Order
	args  *front.WxPayArgs
}

type prepayOrderInput struct {
	tokUsr     *models.User
	orderId    uint
	ip         *string
	chanResult chan *prepayOrderResult
	chanErr    chan error
}

func (hub *OrderHub) PrepayOrder(tokUsr *models.User, orderId uint, cip string) (order *front.Order, args *front.WxPayArgs, err error) {
	chanResult := make(chan *prepayOrderResult)
	chanErr := make(chan error)
	ip := &cip
	in := prepayOrderInput{tokUsr, orderId, ip, chanResult, chanErr}
	hub.chanInput <- in

	select {
	case result := <-chanResult:
		order, args = result.order, result.args
	case err = <-chanErr:
	}
	return
}

func (hub *OrderHub) onPrepayOrder(in *prepayOrderInput) {
	var order *front.Order
	var args *front.WxPayArgs
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		order, args, err = tx.PrepayOrder(in.tokUsr, in.orderId, in.ip)
		return
	})
	if err != nil {
		in.chanErr <- err
	} else {
		in.chanResult <- &prepayOrderResult{order, args}
	}
}
