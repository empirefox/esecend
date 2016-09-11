package hub

import (
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/golang/glog"
)

type ProductHub struct {
	dbs      *dbsrv.DbService
	checkout chan *checkoutInput
}

func NewProductHub(dbs *dbsrv.DbService) *ProductHub {
	return &ProductHub{
		dbs:      dbs,
		checkout: make(chan *checkoutInput, 100),
	}
}

func (hub *ProductHub) Run() {
	for {
		select {
		case c := <-hub.checkout:
			var r *front.Order
			var e error
			switch payload := c.payload.(type) {
			case *front.CheckoutPayload:
				r, e = hub.dbs.CheckoutOrder(c.tokUsr, payload)
			case *front.CheckoutOnePayload:
				r, e = hub.dbs.CheckoutOrderOne(c.tokUsr, payload)
			default:
				glog.Errorf("cannot handle payload: %#v\n", payload)
				e = cerr.Error
			}
			if e != nil {
				c.chanErr <- e
			} else {
				c.chanOrder <- r
			}
		}
	}
}

func (hub *ProductHub) CheckoutOrder(
	tokUsr *models.User, payload *front.CheckoutPayload,
) (order *front.Order, err error) {
	return hub.checkoutOrder(tokUsr, payload)
}

func (hub *ProductHub) CheckoutOrderOne(
	tokUsr *models.User, payload *front.CheckoutOnePayload,
) (order *front.Order, err error) {
	return hub.checkoutOrder(tokUsr, payload)
}

type checkoutInput struct {
	tokUsr    *models.User
	payload   interface{}
	chanOrder chan *front.Order
	chanErr   chan error
}

func (hub *ProductHub) checkoutOrder(tokUsr *models.User, payload interface{}) (order *front.Order, err error) {

	chanOrder := make(chan *front.Order)
	chanErr := make(chan error)
	input := &checkoutInput{tokUsr, payload, chanOrder, chanErr}
	hub.checkout <- input
	select {
	case order = <-chanOrder:
	case err = <-chanErr:
	}
	return
}
