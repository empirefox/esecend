//go:generate reform
package front

import "github.com/empirefox/reform"

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
	TOrderStateHistory
)

//reform:cc_order_item
type OrderItem struct {
	ID         uint    `reform:"id,pk"`
	UserID     uint    `reform:"user_id" json:"-"` // check owner only
	OrderID    uint    `reform:"order_id"`
	ProductID  uint    `reform:"product_id"`
	SkuID      uint    `reform:"sku_id"`
	StoreID    uint    `reform:"store_id"`
	Vpn        VpnType `reform:"vpn"`
	Quantity   uint    `reform:"product_count"`
	Price      uint    `reform:"product_price"`
	CreatedAt  int64   `reform:"created_time"`
	Name       string  `reform:"name"`
	Img        string  `reform:"img"`
	Attrs      string  `reform:"attrs"`
	DeliverFee uint    `reform:"deliver_fee"`

	Store1 uint `reform:"store1"  json:"-"`

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
	ID              uint   `reform:"id,pk"`
	PayAmount       uint   `reform:"pay_amount"`
	WxPaid          uint   `reform:"wx_paid"`
	WxRefund        uint   `reform:"wx_refund"`
	CashPaid        uint   `reform:"cash_paid"`
	CashRefund      uint   `reform:"cash_refund"`
	PayPoints       uint   `reform:"pay_points"` // override other pay types
	PointsPaid      uint   `reform:"points_paid"`
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
	HistoryAt       int64      `reform:"history_at"`

	// auto set by system
	AutoCompleted bool `reform:"auto_completed"`
	AutoEvaled    bool `reform:"auto_evaled"`

	Rebated bool `reform:"rebated" json:"-"`
	User1   uint `reform:"user1"   json:"-"`

	// weixin
	WxPrepayID      string     `reform:"wx_prepay_id"      json:"-"`
	WxTransactionId string     `reform:"wx_transaction_id" json:"-"`
	WxTradeState    TradeState `reform:"wx_trade_state"    json:"-"`
	WxRefundID      string     `reform:"wx_refund_id"      json:"-"`
	WxTradeNo       string     `reform:"wx_trade_no"       json:"-"`

	// Invoice
	InvoiceTo    string `reform:"invoice_to"`
	InvoiceToCom bool   `reform:"invoice_to_com"`

	// OrderAddress
	Contact        string `reform:"contact"`
	Phone          string `reform:"phone"`
	DeliverAddress string `reform:"deliver_addr"`
}

type CheckoutPayloadItem struct {
	SkuID      uint
	Quantity   uint
	GroupBuyID uint
	Attrs      []uint // attr ids sorted

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

type CheckoutOnePayload struct {
	OrderAddress
	Invoice // only for abc
	SkuID   uint
	Attrs   []uint // attr ids sorted
	Remark  string
}

type OrderChangeStatePayload struct {
	ID    uint
	State OrderState
}

type OrderPrepayPayload struct {
	OrderID uint
}

type OrderPrepayResponse struct {
	Order     *Order
	WxPayArgs *WxPayArgs
}

type OrderPayPayload struct {
	Key      string
	OrderID  uint
	Amount   uint
	IsPoints bool
}

type OrdersResponse struct {
	Orders []reform.Struct // Order
	Items  []reform.Struct // OrderItem
}
