//go:generate reform
package front

type TradeState int

const (
	UNKNOWN TradeState = iota
	NOTPAY
	SUCCESS
	REFUND
	CLOSED
	REVOKED
	USERPAYING
	PAYERROR
)

type OrderState int

const (
	TOrderStateUnknown OrderState = iota
	TOrderStateNopay
	TOrderStatePrepaid
	TOrderStatePaid
	TOrderStateCanceled
	TOrderStatePicking
	TOrderStateDelivered
	TOrderStateReturnStarted
	TOrderStateReturning
	TOrderStateReturned
	TOrderStateRejecting
	TOrderStateRejectBack
	TOrderStateRejectRefound
	TOrderStateCompleted
	TOrderStateEvalStarted
	TOrderStateEvaled
)

//reform:cc_order_item
type OrderItem struct {
	ID         uint   `reform:"id,pk"`
	UserID     uint   `reform:"user_id" json:"-"` // check owner only
	OrderID    uint   `reform:"order_id"`
	ProductID  uint   `reform:"product_id"`
	SkuID      uint   `reform:"sku_id"`
	Quantity   uint   `reform:"product_count"`
	Price      uint   `reform:"product_price"`
	CreatedAt  int64  `reform:"created_time"`
	Name       string `reform:"name"`
	Img        string `reform:"img"`
	Attrs      string `reform:"attrs"`
	DeliverFee uint   `reform:"deliver_fee"`

	// EvalItem
	Eval        string `reform:"comment_content"`
	EvalAt      int64  `reform:"comment_time"` // gen by server
	EvalName    string `reform:"user_name"`    // gen by server
	RateStar    uint   `reform:"rate_star"`
	RateFit     uint   `reform:"rate_fit"`
	RateServe   uint   `reform:"rate_serve"`
	RateDeliver uint   `reform:"rate_deliver"`
}

//reform:cc_order
type Order struct {
	//	CashPaidID      uint   `reform:"cash_paid_id"`
	//	PointsPaidID    uint   `reform:"points_paid_id"`
	ID              uint   `reform:"id,pk"`
	PayAmount       uint   `reform:"pay_amount"`
	WxPaid          uint   `reform:"wx_paid"`
	CashPaid        uint   `reform:"cash_paid"`
	PointsPaid      uint   `reform:"points_paid"`
	WxRefund        uint   `reform:"wx_refund"`
	CashRefund      uint   `reform:"cash_refund"`
	PointsRefund    uint   `reform:"points_refund"`
	AbandonedReason string `reform:"abandoned_reason"`
	Remark          string `reform:"remark"`
	UserID          uint   `reform:"user_id" json:"-"` // check owner only
	Items           []*OrderItem

	IsDeliverPay bool   `reform:"is_deliver_pay"`
	DeliverFee   uint   `reform:"deliver_fee"`
	DeliverCom   string `reform:"deliver_com"`
	DeliverNo    string `reform:"deliver_no"`

	// study http://help.vipshop.com/themelist.php?type=detail&id=330
	State           OrderState `reform:"state"`
	CreatedAt       int64      `reform:"created_at"`
	CanceledAt      int64      `reform:"canceled_at"`
	PaidAt          int64      `reform:"paid_at"`
	PrepaidAt       int64      `reform:"prepaid_at"`
	PaidCanceledAt  int64      `reform:"paid_canceled_at"` // auto refound
	PickingAt       int64      `reform:"picking_at"`       // min 5 minute from paid, no cancel
	DeliveredAt     int64      `reform:"delivered_at"`
	ReturnStaredtAt int64      `reform:"return_started_at"`
	ReturnEnsuredAt int64      `reform:"return_ensured_at"`
	ReturnedAt      int64      `reform:"returned_at"` // auto refound
	RejectingAt     int64      `reform:"rejecting_at"`
	RejectBackAt    int64      `reform:"reject_back_at"`    // goods back
	RejectRefoundAt int64      `reform:"reject_refound_at"` // auto refound
	CompletedAt     int64      `reform:"completed_at"`
	EvalStartedAt   int64      `reform:"eval_started_at"`
	EvalAt          int64      `reform:"eval_at"`

	// weixin
	TransactionId string     `reform:"transaction_id"`
	TradeState    TradeState `reform:"trade_state"`
	WxRefundID    string     `reform:"refund_id"`

	// Invoice
	InvoiceTo    string `reform:"invoice_to"`
	InvoiceToCom bool   `reform:"invoice_to_com"`

	// OrderAddress
	Contact        string `reform:"contact"`
	Phone          string `reform:"phone"`
	DeliverAddress string `reform:"deliver_addr"`
}

type CheckoutPayloadItem struct {
	SkuID         uint
	Quantity      uint
	SkuPrice      uint // TODO
	GroupBuyID    uint
	GroupBuyPrice uint   // TODO
	Attrs         []uint // attr ids sorted

	// ids queried from db
	AttrIds []uint `json:"-"`

	// values from db after validate
	AttrValues string `json:"-"`
}

type Invoice struct {
	InvoiceTo    string
	InvoiceToCom bool
}

type OrderAddress struct {
	Contact        string
	Phone          string
	DeliverAddress string
}

type CheckoutPayload struct {
	OrderAddress
	Invoice
	Items        []CheckoutPayloadItem
	Remark       string
	Total        uint // final amount to pay, used to validate
	DeliverFee   uint
	IsDeliverPay bool
}

type OrderChangeStatePayload struct {
	ID    uint
	State OrderState
}

type OrderPrepayPayload struct {
	OrderID uint
	Amount  uint

	Cash   uint
	Wx     uint
	Points uint // just points not cents

	Ip string
}

type OrderPrepayResponse struct {
	Order     *Order
	WxPayArgs *WxPayArgs
}

type OrderPayPayload struct {
	Key     string
	Amount  uint
	OrderID uint

	Cash   uint
	Points uint // just points not cents
}
