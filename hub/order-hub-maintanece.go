package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/golang/glog"
)

func (hub *OrderHub) onMaintain() {
	err := hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.OrdersMaintain()
		return
	})
	if err != nil {
		glog.Errorln(err)
	}

	err = hub.dbs.InTx(func(tx *dbsrv.DbService) (err error) {
		err = tx.RebateMaintain()
		return
	})
	if err != nil {
		glog.Errorln(err)
	}
}
