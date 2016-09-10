package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/golang/glog"
)

type OrderHub struct {
	dbs       *dbsrv.DbService
	chanInput chan interface{}
}

func NewOrderHub(dbs *dbsrv.DbService) *OrderHub {
	return &OrderHub{
		dbs:       dbs,
		chanInput: make(chan interface{}, 100),
	}
}

func (hub *OrderHub) Run() {
	for input := range hub.chanInput {
		switch in := input.(type) {
		case *prepayOrderInput:
			hub.onPrepayOrder(in)
		case *payOrderInput:
			hub.onPayOrder(in)
		case *orderPaidStateInput:
			hub.onOrderPaidState(in)
		case *orderOnWxPayNotifyInput:
			hub.onWxPayNotify(in)
		case *orderChangeStateInput:
			hub.onOrderChangeState(in)
		case *orderMgrStateInput:
			hub.onMgrOrderState(in)
		case *orderEvalInput:
			hub.onEvalSave(in)
		default:
			glog.Errorf("cannot handle input: %#v\n", in)
		}
	}
}
