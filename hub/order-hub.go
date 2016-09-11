package hub

import (
	"time"

	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/golang/glog"
)

type OrderHub struct {
	config    *config.Order
	dbs       *dbsrv.DbService
	chanInput chan interface{}
}

func NewOrderHub(config *config.Config, dbs *dbsrv.DbService) *OrderHub {
	return &OrderHub{
		config:    &config.Order,
		dbs:       dbs,
		chanInput: make(chan interface{}, 100),
	}
}

func (hub *OrderHub) Run() {
	ticker := time.NewTicker(time.Duration(hub.config.MaintaneTimeMinute) * time.Minute)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case input := <-hub.chanInput:
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

		case <-ticker.C:
			hub.onOrderMaintanece()
		}
	}
}
