package front

import (
	"fmt"

	"github.com/empirefox/esecend/fsm"
)

func (o *Order) TrackingNumber() string {
	return fmt.Sprintf("%d-%d", o.CreatedAt, o.ID)
}

func (o *Order) CurrentState() fsm.State { return fsm.State(o.State) }
func (o *Order) SetState(s fsm.State)    { o.State = OrderState(s) }
