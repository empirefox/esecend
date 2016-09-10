package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/golang/glog"
)

func (hub *OrderHub) onOrderMaintanece() {
	err := hub.dbs.IsOrderCompleted(func(tx *dbsrv.DbService) (err error) {
		err = tx.OrdersMaintanence()
		return
	})
	if err != nil {
		glog.Errorln(err)
	}
}
