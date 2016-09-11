package hub

import (
	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
)

type orderMgrStateInput struct {
	order   *front.Order
	claims  *admin.Claims
	chanErr chan error
}

func (hub *OrderHub) MgrOrderState(
	order *front.Order, claims *admin.Claims,
) (err error) {
	chanErr := make(chan error)
	in := &orderMgrStateInput{order, claims, chanErr}
	hub.chanInput <- in

	err = <-chanErr
	return
}

func (hub *OrderHub) onMgrOrderState(in *orderMgrStateInput) {
	in.chanErr <- hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.MgrOrderState(in.order, in.claims)
		return
	})
}
