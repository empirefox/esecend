package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/golang/glog"
)

func (hub *OrderHub) onOrderMaintanece() {
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.OrdersMaintanence()
		return
	})
	if err != nil {
		glog.Errorln(err)
	}
}
