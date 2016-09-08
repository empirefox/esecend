package hub

import "github.com/empirefox/esecend/db-service"

type OrderHub struct {
	dbs       *dbsrv.DbService
	chanInput chan interface{}
}

func (hub *OrderHub) Run() {
	for input := range hub.chanInput {
		switch in := input.(type) {
		case *prepayOrderInput:
			hub.onPrepayOrder(tx, in)
		case *payOrderInput:
			hub.onPayOrder(tx, in)
		}
	}
}
