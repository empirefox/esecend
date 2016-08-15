package dbsrv

import (
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/fsm"
)

var orderRules = fsm.CreateRuleset(
	newTransition(front.TOrderStateNopay, front.TOrderStatePrepaid),
	newTransition(front.TOrderStateNopay, front.TOrderStatePaid),
	newTransition(front.TOrderStateNopay, front.TOrderStateCanceled),

	newTransition(front.TOrderStatePrepaid, front.TOrderStatePaid),
	newTransition(front.TOrderStatePrepaid, front.TOrderStateCanceled),

	newTransition(front.TOrderStatePaid, front.TOrderStatePicking), // admin
	newTransition(front.TOrderStatePaid, front.TOrderStateCanceled),

	newTransition(front.TOrderStatePicking, front.TOrderStateDelivered), // admin
	newTransition(front.TOrderStatePicking, front.TOrderStateCanceled),

	newTransition(front.TOrderStateDelivered, front.TOrderStateCompleted),   // user+system
	newTransition(front.TOrderStateDelivered, front.TOrderStateEvalStarted), // standalone
	newTransition(front.TOrderStateDelivered, front.TOrderStateEvaled),      // standalone
	newTransition(front.TOrderStateDelivered, front.TOrderStateRejecting),   // admin
	newTransition(front.TOrderStateDelivered, front.TOrderStateReturnStarted),

	newTransition(front.TOrderStateCompleted, front.TOrderStateEvalStarted), // standalone
	newTransition(front.TOrderStateCompleted, front.TOrderStateEvaled),      // standalone

	newTransition(front.TOrderStateEvaled, front.TOrderStateEvaled),      // standalone
	newTransition(front.TOrderStateEvalStarted, front.TOrderStateEvaled), // standalone

	newTransition(front.TOrderStateRejecting, front.TOrderStateRejectBack),     // admin
	newTransition(front.TOrderStateRejectBack, front.TOrderStateRejectRefound), // admin

	newTransition(front.TOrderStateReturnStarted, front.TOrderStateReturning), // admin
	newTransition(front.TOrderStateReturning, front.TOrderStateReturned),      // admin
)

func newTransition(from, to front.OrderState) fsm.T {
	return fsm.T{fsm.State(from), fsm.State(to)}
}

func PermitOrderState(o *front.Order, s front.OrderState) error {
	return orderRules.Permitted(o, fsm.State(s))
}
