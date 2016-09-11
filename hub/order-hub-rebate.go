package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

type userVipRebateInput struct {
	tokUsr  *models.User
	payload *front.VipRebateRequest
	chanErr chan<- error
}

func (hub *OrderHub) UserVipRebate(tokUsr *models.User, payload *front.VipRebateRequest) (err error) {
	chanErr := make(chan error)
	in := &userVipRebateInput{tokUsr, payload, chanErr}
	hub.chanInput <- in

	err = <-chanErr
	return
}

func (hub *OrderHub) onUserVipRebate(in *userVipRebateInput) {
	in.chanErr <- hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.UserVipRebate(in.tokUsr, in.payload)
		return
	})
}
