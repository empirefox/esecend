package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type orderItemTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *orderItemTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_order_item").
func (v *orderItemTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *orderItemTable) Columns() []string {
	return []string{"id", "user_id", "order_id", "product_id", "sku_id", "product_count", "product_price", "created_time", "name", "img", "attrs", "deliver_fee", "comment_content", "comment_time", "user_name", "rate_star", "rate_fit", "rate_serve", "rate_deliver"}
}

// NewStruct makes a new struct for that view or table.
func (v *orderItemTable) NewStruct() reform.Struct {
	return new(OrderItem)
}

// NewRecord makes a new record for that table.
func (v *orderItemTable) NewRecord() reform.Record {
	return new(OrderItem)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *orderItemTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// OrderItemTable represents cc_order_item view or table in SQL database.
var OrderItemTable = &orderItemTable{
	s: parse.StructInfo{Type: "OrderItem", SQLSchema: "", SQLName: "cc_order_item", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "UserID", PKType: "", Column: "user_id"}, {Name: "OrderID", PKType: "", Column: "order_id"}, {Name: "ProductID", PKType: "", Column: "product_id"}, {Name: "SkuID", PKType: "", Column: "sku_id"}, {Name: "Quantity", PKType: "", Column: "product_count"}, {Name: "Price", PKType: "", Column: "product_price"}, {Name: "CreatedAt", PKType: "", Column: "created_time"}, {Name: "Name", PKType: "", Column: "name"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Attrs", PKType: "", Column: "attrs"}, {Name: "DeliverFee", PKType: "", Column: "deliver_fee"}, {Name: "Eval", PKType: "", Column: "comment_content"}, {Name: "EvalAt", PKType: "", Column: "comment_time"}, {Name: "EvalName", PKType: "", Column: "user_name"}, {Name: "RateStar", PKType: "", Column: "rate_star"}, {Name: "RateFit", PKType: "", Column: "rate_fit"}, {Name: "RateServe", PKType: "", Column: "rate_serve"}, {Name: "RateDeliver", PKType: "", Column: "rate_deliver"}}, PKFieldIndex: 0},
	z: new(OrderItem).Values(),
}

// String returns a string representation of this struct or record.
func (s OrderItem) String() string {
	res := make([]string, 19)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "OrderID: " + reform.Inspect(s.OrderID, true)
	res[3] = "ProductID: " + reform.Inspect(s.ProductID, true)
	res[4] = "SkuID: " + reform.Inspect(s.SkuID, true)
	res[5] = "Quantity: " + reform.Inspect(s.Quantity, true)
	res[6] = "Price: " + reform.Inspect(s.Price, true)
	res[7] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[8] = "Name: " + reform.Inspect(s.Name, true)
	res[9] = "Img: " + reform.Inspect(s.Img, true)
	res[10] = "Attrs: " + reform.Inspect(s.Attrs, true)
	res[11] = "DeliverFee: " + reform.Inspect(s.DeliverFee, true)
	res[12] = "Eval: " + reform.Inspect(s.Eval, true)
	res[13] = "EvalAt: " + reform.Inspect(s.EvalAt, true)
	res[14] = "EvalName: " + reform.Inspect(s.EvalName, true)
	res[15] = "RateStar: " + reform.Inspect(s.RateStar, true)
	res[16] = "RateFit: " + reform.Inspect(s.RateFit, true)
	res[17] = "RateServe: " + reform.Inspect(s.RateServe, true)
	res[18] = "RateDeliver: " + reform.Inspect(s.RateDeliver, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *OrderItem) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.OrderID,
		s.ProductID,
		s.SkuID,
		s.Quantity,
		s.Price,
		s.CreatedAt,
		s.Name,
		s.Img,
		s.Attrs,
		s.DeliverFee,
		s.Eval,
		s.EvalAt,
		s.EvalName,
		s.RateStar,
		s.RateFit,
		s.RateServe,
		s.RateDeliver,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *OrderItem) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.OrderID,
		&s.ProductID,
		&s.SkuID,
		&s.Quantity,
		&s.Price,
		&s.CreatedAt,
		&s.Name,
		&s.Img,
		&s.Attrs,
		&s.DeliverFee,
		&s.Eval,
		&s.EvalAt,
		&s.EvalName,
		&s.RateStar,
		&s.RateFit,
		&s.RateServe,
		&s.RateDeliver,
	}
}

// View returns View object for that struct.
func (s *OrderItem) View() reform.View {
	return OrderItemTable
}

// Table returns Table object for that record.
func (s *OrderItem) Table() reform.Table {
	return OrderItemTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *OrderItem) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *OrderItem) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *OrderItem) HasPK() bool {
	return s.ID != OrderItemTable.z[OrderItemTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *OrderItem) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = OrderItemTable
	_ reform.Struct = new(OrderItem)
	_ reform.Table  = OrderItemTable
	_ reform.Record = new(OrderItem)
	_ fmt.Stringer  = new(OrderItem)
)

type orderTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *orderTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_order").
func (v *orderTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *orderTable) Columns() []string {
	return []string{"id", "pay_amount", "wx_paid", "cash_paid", "points_paid", "wx_refund", "cash_refund", "points_refund", "abandoned_reason", "remark", "user_id", "is_deliver_pay", "deliver_fee", "deliver_com", "deliver_no", "state", "created_at", "canceled_at", "paid_at", "prepaid_at", "paid_canceled_at", "picking_at", "delivered_at", "return_started_at", "return_ensured_at", "returned_at", "rejecting_at", "reject_back_at", "reject_refound_at", "completed_at", "eval_started_at", "eval_at", "transaction_id", "trade_state", "refund_id", "invoice_to", "invoice_to_com", "contact", "phone", "deliver_addr"}
}

// NewStruct makes a new struct for that view or table.
func (v *orderTable) NewStruct() reform.Struct {
	return new(Order)
}

// NewRecord makes a new record for that table.
func (v *orderTable) NewRecord() reform.Record {
	return new(Order)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *orderTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// OrderTable represents cc_order view or table in SQL database.
var OrderTable = &orderTable{
	s: parse.StructInfo{Type: "Order", SQLSchema: "", SQLName: "cc_order", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "PayAmount", PKType: "", Column: "pay_amount"}, {Name: "WxPaid", PKType: "", Column: "wx_paid"}, {Name: "CashPaid", PKType: "", Column: "cash_paid"}, {Name: "PointsPaid", PKType: "", Column: "points_paid"}, {Name: "WxRefund", PKType: "", Column: "wx_refund"}, {Name: "CashRefund", PKType: "", Column: "cash_refund"}, {Name: "PointsRefund", PKType: "", Column: "points_refund"}, {Name: "AbandonedReason", PKType: "", Column: "abandoned_reason"}, {Name: "Remark", PKType: "", Column: "remark"}, {Name: "UserID", PKType: "", Column: "user_id"}, {Name: "IsDeliverPay", PKType: "", Column: "is_deliver_pay"}, {Name: "DeliverFee", PKType: "", Column: "deliver_fee"}, {Name: "DeliverCom", PKType: "", Column: "deliver_com"}, {Name: "DeliverNo", PKType: "", Column: "deliver_no"}, {Name: "State", PKType: "", Column: "state"}, {Name: "CreatedAt", PKType: "", Column: "created_at"}, {Name: "CanceledAt", PKType: "", Column: "canceled_at"}, {Name: "PaidAt", PKType: "", Column: "paid_at"}, {Name: "PrepaidAt", PKType: "", Column: "prepaid_at"}, {Name: "PaidCanceledAt", PKType: "", Column: "paid_canceled_at"}, {Name: "PickingAt", PKType: "", Column: "picking_at"}, {Name: "DeliveredAt", PKType: "", Column: "delivered_at"}, {Name: "ReturnStaredtAt", PKType: "", Column: "return_started_at"}, {Name: "ReturnEnsuredAt", PKType: "", Column: "return_ensured_at"}, {Name: "ReturnedAt", PKType: "", Column: "returned_at"}, {Name: "RejectingAt", PKType: "", Column: "rejecting_at"}, {Name: "RejectBackAt", PKType: "", Column: "reject_back_at"}, {Name: "RejectRefoundAt", PKType: "", Column: "reject_refound_at"}, {Name: "CompletedAt", PKType: "", Column: "completed_at"}, {Name: "EvalStartedAt", PKType: "", Column: "eval_started_at"}, {Name: "EvalAt", PKType: "", Column: "eval_at"}, {Name: "TransactionId", PKType: "", Column: "transaction_id"}, {Name: "TradeState", PKType: "", Column: "trade_state"}, {Name: "WxRefundID", PKType: "", Column: "refund_id"}, {Name: "InvoiceTo", PKType: "", Column: "invoice_to"}, {Name: "InvoiceToCom", PKType: "", Column: "invoice_to_com"}, {Name: "Contact", PKType: "", Column: "contact"}, {Name: "Phone", PKType: "", Column: "phone"}, {Name: "DeliverAddress", PKType: "", Column: "deliver_addr"}}, PKFieldIndex: 0},
	z: new(Order).Values(),
}

// String returns a string representation of this struct or record.
func (s Order) String() string {
	res := make([]string, 40)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "PayAmount: " + reform.Inspect(s.PayAmount, true)
	res[2] = "WxPaid: " + reform.Inspect(s.WxPaid, true)
	res[3] = "CashPaid: " + reform.Inspect(s.CashPaid, true)
	res[4] = "PointsPaid: " + reform.Inspect(s.PointsPaid, true)
	res[5] = "WxRefund: " + reform.Inspect(s.WxRefund, true)
	res[6] = "CashRefund: " + reform.Inspect(s.CashRefund, true)
	res[7] = "PointsRefund: " + reform.Inspect(s.PointsRefund, true)
	res[8] = "AbandonedReason: " + reform.Inspect(s.AbandonedReason, true)
	res[9] = "Remark: " + reform.Inspect(s.Remark, true)
	res[10] = "UserID: " + reform.Inspect(s.UserID, true)
	res[11] = "IsDeliverPay: " + reform.Inspect(s.IsDeliverPay, true)
	res[12] = "DeliverFee: " + reform.Inspect(s.DeliverFee, true)
	res[13] = "DeliverCom: " + reform.Inspect(s.DeliverCom, true)
	res[14] = "DeliverNo: " + reform.Inspect(s.DeliverNo, true)
	res[15] = "State: " + reform.Inspect(s.State, true)
	res[16] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[17] = "CanceledAt: " + reform.Inspect(s.CanceledAt, true)
	res[18] = "PaidAt: " + reform.Inspect(s.PaidAt, true)
	res[19] = "PrepaidAt: " + reform.Inspect(s.PrepaidAt, true)
	res[20] = "PaidCanceledAt: " + reform.Inspect(s.PaidCanceledAt, true)
	res[21] = "PickingAt: " + reform.Inspect(s.PickingAt, true)
	res[22] = "DeliveredAt: " + reform.Inspect(s.DeliveredAt, true)
	res[23] = "ReturnStaredtAt: " + reform.Inspect(s.ReturnStaredtAt, true)
	res[24] = "ReturnEnsuredAt: " + reform.Inspect(s.ReturnEnsuredAt, true)
	res[25] = "ReturnedAt: " + reform.Inspect(s.ReturnedAt, true)
	res[26] = "RejectingAt: " + reform.Inspect(s.RejectingAt, true)
	res[27] = "RejectBackAt: " + reform.Inspect(s.RejectBackAt, true)
	res[28] = "RejectRefoundAt: " + reform.Inspect(s.RejectRefoundAt, true)
	res[29] = "CompletedAt: " + reform.Inspect(s.CompletedAt, true)
	res[30] = "EvalStartedAt: " + reform.Inspect(s.EvalStartedAt, true)
	res[31] = "EvalAt: " + reform.Inspect(s.EvalAt, true)
	res[32] = "TransactionId: " + reform.Inspect(s.TransactionId, true)
	res[33] = "TradeState: " + reform.Inspect(s.TradeState, true)
	res[34] = "WxRefundID: " + reform.Inspect(s.WxRefundID, true)
	res[35] = "InvoiceTo: " + reform.Inspect(s.InvoiceTo, true)
	res[36] = "InvoiceToCom: " + reform.Inspect(s.InvoiceToCom, true)
	res[37] = "Contact: " + reform.Inspect(s.Contact, true)
	res[38] = "Phone: " + reform.Inspect(s.Phone, true)
	res[39] = "DeliverAddress: " + reform.Inspect(s.DeliverAddress, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Order) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.PayAmount,
		s.WxPaid,
		s.CashPaid,
		s.PointsPaid,
		s.WxRefund,
		s.CashRefund,
		s.PointsRefund,
		s.AbandonedReason,
		s.Remark,
		s.UserID,
		s.IsDeliverPay,
		s.DeliverFee,
		s.DeliverCom,
		s.DeliverNo,
		s.State,
		s.CreatedAt,
		s.CanceledAt,
		s.PaidAt,
		s.PrepaidAt,
		s.PaidCanceledAt,
		s.PickingAt,
		s.DeliveredAt,
		s.ReturnStaredtAt,
		s.ReturnEnsuredAt,
		s.ReturnedAt,
		s.RejectingAt,
		s.RejectBackAt,
		s.RejectRefoundAt,
		s.CompletedAt,
		s.EvalStartedAt,
		s.EvalAt,
		s.TransactionId,
		s.TradeState,
		s.WxRefundID,
		s.InvoiceTo,
		s.InvoiceToCom,
		s.Contact,
		s.Phone,
		s.DeliverAddress,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Order) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.PayAmount,
		&s.WxPaid,
		&s.CashPaid,
		&s.PointsPaid,
		&s.WxRefund,
		&s.CashRefund,
		&s.PointsRefund,
		&s.AbandonedReason,
		&s.Remark,
		&s.UserID,
		&s.IsDeliverPay,
		&s.DeliverFee,
		&s.DeliverCom,
		&s.DeliverNo,
		&s.State,
		&s.CreatedAt,
		&s.CanceledAt,
		&s.PaidAt,
		&s.PrepaidAt,
		&s.PaidCanceledAt,
		&s.PickingAt,
		&s.DeliveredAt,
		&s.ReturnStaredtAt,
		&s.ReturnEnsuredAt,
		&s.ReturnedAt,
		&s.RejectingAt,
		&s.RejectBackAt,
		&s.RejectRefoundAt,
		&s.CompletedAt,
		&s.EvalStartedAt,
		&s.EvalAt,
		&s.TransactionId,
		&s.TradeState,
		&s.WxRefundID,
		&s.InvoiceTo,
		&s.InvoiceToCom,
		&s.Contact,
		&s.Phone,
		&s.DeliverAddress,
	}
}

// View returns View object for that struct.
func (s *Order) View() reform.View {
	return OrderTable
}

// Table returns Table object for that record.
func (s *Order) Table() reform.Table {
	return OrderTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Order) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Order) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Order) HasPK() bool {
	return s.ID != OrderTable.z[OrderTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Order) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = OrderTable
	_ reform.Struct = new(Order)
	_ reform.Table  = OrderTable
	_ reform.Record = new(Order)
	_ fmt.Stringer  = new(Order)
)

func init() {
	parse.AssertUpToDate(&OrderItemTable.s, new(OrderItem))
	OrderItemTable.ViewBase = reform.NewViewBase(&OrderItemTable.s)
	parse.AssertUpToDate(&OrderTable.s, new(Order))
	OrderTable.ViewBase = reform.NewViewBase(&OrderTable.s)
}
