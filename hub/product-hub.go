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

func NewProductHub() *ProductHub {
	return &ProductHub{
		checkout: make(chan *checkoutInput, 100),
	}
}

func (hub *ProductHub) Run() {
	for {
		select {
		case c := <-checkout:
			r, e := dbs.CheckoutOrder(c.tokUsr, c.payload)
			if e != nil {
				c.err <- e
			} else {
				c.result <- r
			}
		}
	}
}

type checkoutInput struct {
	tokUsr  *models.User
	payload *front.CheckoutPayload
	result  <-chan *front.Order
	err     <-chan error
}

func (hub *ProductHub) CheckoutOrder(tokUsr *models.User, payload *front.CheckoutPayload) (*front.Order, error) {
	result := make(chan *front.Order)
	err := make(chan error)
	input := &checkoutInput{tokUsr, payload, result, err}
	hub.checkout <- input
	select {
	case r := <-result:
		return r, nil
	case e := <-err:
		return nil, e
	}
}
