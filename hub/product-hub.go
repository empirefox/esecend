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
		case c := <-checkout:
			r, e := dbs.CheckoutOrder(c.tokUsr, c.payload)
			if e != nil {
				c.chanErr <- e
			} else {
				c.chanOrder <- r
			}
		}
	}
}

type checkoutInput struct {
	tokUsr    *models.User
	payload   *front.CheckoutPayload
	chanOrder <-chan *front.Order
	chanErr   <-chan error
}

func (hub *ProductHub) CheckoutOrder(tokUsr *models.User, payload *front.CheckoutPayload) (order *front.Order, err error) {
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
