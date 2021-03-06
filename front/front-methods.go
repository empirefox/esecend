package front

import (
	"fmt"

	"github.com/empirefox/esecend/fsm"
)

func (o *Order) TrackingNumber() string {
	return fmt.Sprintf("%d-%d", o.CreatedAt, o.ID)
}

func (o *Order) WxOutTradeNo() string {
	return fmt.Sprintf("%d-%d-%s", o.PrepaidAt, o.ID, o.WxTradeNo)
}

func (o *Order) CurrentState() fsm.State { return fsm.State(o.State) }
func (o *Order) SetState(s fsm.State)    { o.State = OrderState(s) }

func (c *UserCash) TrackingNumber() string {
	return fmt.Sprintf("%d-%d", c.CreatedAt, c.ID)
}
