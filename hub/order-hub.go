package hub

import (
	"fmt"
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

func (hub *OrderHub) run(input interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
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

	case *userVipRebateInput:
		hub.onUserVipRebate(in)

	case *userWithdrawInput:
		hub.onUserWithdraw(in)

	default:
		glog.Errorf("cannot handle input: %#v\n", in)
	}
}

func (hub *OrderHub) Run() {
	hub.onMaintain()
	ticker := time.NewTicker(time.Duration(hub.config.MaintainTimeMinute) * time.Minute)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case input := <-hub.chanInput:
			hub.run(input)

		case <-ticker.C:
			hub.onMaintain()
		}
	}
}
