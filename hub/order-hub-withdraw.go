package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type userWithdrawInput struct {
	tokUsr  *models.User
	payload *front.WithdrawPayload
	chanErr <-chan error
}

func (hub *OrderHub) UserWithdraw(tokUsr *models.User, payload *front.WithdrawPayload) (err error) {
	chanErr := make(chan error)
	in := &userWithdrawInput{tokUsr, payload, chanErr}
	hub.chanInput <- in

	err = <-chanErr
	return
}

func (hub *OrderHub) onUserWithdraw(in *userWithdrawInput) {
	in.chanErr <- hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.UserWithdraw(in.tokUsr, in.payload)
		return
	})
}
