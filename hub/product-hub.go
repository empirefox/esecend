package hub

import (
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
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
			if c.isOne {
				r, e = dbs.CheckoutOrderOne(c.tokUsr, c.payload)
			} else {
				r, e = dbs.CheckoutOrder(c.tokUsr, c.payload)
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
	return hub.checkoutOrder(tokUsr, payload, false)
}

func (hub *ProductHub) CheckoutOrderOne(
	tokUsr *models.User, payload *front.CheckoutPayload,
) (order *front.Order, err error) {
	return hub.checkoutOrder(tokUsr, payload, true)
}

type checkoutInput struct {
	tokUsr    *models.User
	payload   *front.CheckoutPayload
	isOne     bool
	chanOrder <-chan *front.Order
	chanErr   <-chan error
}

func (hub *ProductHub) checkoutOrder(
	tokUsr *models.User, payload *front.CheckoutPayload, isOne bool,
) (order *front.Order, err error) {

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
