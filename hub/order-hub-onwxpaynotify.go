package hub

import "github.com/empirefox/esecend/db-service"

type orderOnWxPayNotifyInput struct {
	src     map[string]string
	orderId uint
	chanErr <-chan error
}

func (hub *OrderHub) OnWxPayNotify(src map[string]string, orderId uint) (err error) {
	chanErr := make(chan error)
	in := &orderOnWxPayNotifyInput{src, orderId, chanErr}
	hub.chanInput <- in

	err = <-chanErr
	return
}

func (hub *OrderHub) onWxPayNotify(in *orderOnWxPayNotifyInput) {
	in.chanErr <- hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.OnWxPayNotify(in.src, in.orderId)
		return
	})
}
