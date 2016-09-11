package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type userWithdrawInput struct {
	tokUsr   *models.User
	payload  *front.WithdrawPayload
	chanCash chan *front.UserCash
	chanErr  chan error
}

func (hub *OrderHub) UserWithdraw(
	tokUsr *models.User, payload *front.WithdrawPayload,
) (cash *front.UserCash, err error) {

	chanCash := make(chan *front.UserCash)
	chanErr := make(chan error)
	in := &userWithdrawInput{tokUsr, payload, chanCash, chanErr}
	hub.chanInput <- in

	select {
	case cash = <-chanCash:
	case err = <-chanErr:
	}
	return
}

func (hub *OrderHub) onUserWithdraw(in *userWithdrawInput) {
	var cash *front.UserCash
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		cash, err = tx.UserWithdraw(in.tokUsr, in.payload)
		return
	})
	if err != nil {
		in.chanErr <- err
	} else {
		in.chanCash <- cash
	}
}
