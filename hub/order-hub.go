package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
)

type OrderHub struct {
	dbs       *dbsrv.DbService
	chanInput chan interface{}
}

func (hub *OrderHub) Run() {
	for input := range hub.chanInput {
		switch in := input.(type) {
		case *prepayOrderInput:
			hub.onPrepayOrder(in)
		}
	}
}

type prepayOrderInput struct {
	*dbsrv.PrepayOrderInput
	chanArgs <-chan *front.WxPayArgs
	chanErr  <-chan error
}

func (hub *OrderHub) PrepayOrder(input *dbsrv.PrepayOrderInput) (args *front.WxPayArgs, err error) {
	chanArgs := make(chan *front.WxPayArgs)
	chanErr := make(chan error)
	in := prepayOrderInput{
		PrepayOrderInput: input,
		chanArgs:         chanArgs,
		chanErr:          chanErr,
	}
	hub.chanInput <- in

	select {
	case args = <-chanArgs:
	case err = <-chanErr:
	}
	return
}

func (hub *OrderHub) onPrepayOrder(in *prepayOrderInput) {
	args, err := hub.dbs.PrepayOrder(in.PrepayOrderInput)
	if err != nil {
		in.chanErr <- err
	} else {
		in.chanArgs <- args
	}
}
