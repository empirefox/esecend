package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type capitalFlowTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *capitalFlowTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_member_account_log").
func (v *capitalFlowTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *capitalFlowTable) Columns() []string {
	return []string{"id", "user_id", "create_time", "log_type", "reason", "amount", "balance", "order_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *capitalFlowTable) NewStruct() reform.Struct {
	return new(CapitalFlow)
}

// NewRecord makes a new record for that table.
func (v *capitalFlowTable) NewRecord() reform.Record {
	return new(CapitalFlow)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *capitalFlowTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// CapitalFlowTable represents cc_member_account_log view or table in SQL database.
var CapitalFlowTable = &capitalFlowTable{
	s: parse.StructInfo{Type: "CapitalFlow", SQLSchema: "", SQLName: "cc_member_account_log", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "UserID", PKType: "", Column: "user_id"}, {Name: "CreatedAt", PKType: "", Column: "create_time"}, {Name: "Type", PKType: "", Column: "log_type"}, {Name: "Reason", PKType: "", Column: "reason"}, {Name: "Amount", PKType: "", Column: "amount"}, {Name: "Balance", PKType: "", Column: "balance"}, {Name: "OrderID", PKType: "", Column: "order_id"}}, PKFieldIndex: 0},
	z: new(CapitalFlow).Values(),
}

// String returns a string representation of this struct or record.
func (s CapitalFlow) String() string {
	res := make([]string, 8)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[3] = "Type: " + reform.Inspect(s.Type, true)
	res[4] = "Reason: " + reform.Inspect(s.Reason, true)
	res[5] = "Amount: " + reform.Inspect(s.Amount, true)
	res[6] = "Balance: " + reform.Inspect(s.Balance, true)
	res[7] = "OrderID: " + reform.Inspect(s.OrderID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *CapitalFlow) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.CreatedAt,
		s.Type,
		s.Reason,
		s.Amount,
		s.Balance,
		s.OrderID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *CapitalFlow) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.CreatedAt,
		&s.Type,
		&s.Reason,
		&s.Amount,
		&s.Balance,
		&s.OrderID,
	}
}

// View returns View object for that struct.
func (s *CapitalFlow) View() reform.View {
	return CapitalFlowTable
}

// Table returns Table object for that record.
func (s *CapitalFlow) Table() reform.Table {
	return CapitalFlowTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *CapitalFlow) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *CapitalFlow) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *CapitalFlow) HasPK() bool {
	return s.ID != CapitalFlowTable.z[CapitalFlowTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *CapitalFlow) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = CapitalFlowTable
	_ reform.Struct = new(CapitalFlow)
	_ reform.Table  = CapitalFlowTable
	_ reform.Record = new(CapitalFlow)
	_ fmt.Stringer  = new(CapitalFlow)
)

type pointsItemTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *pointsItemTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_member_credit_log").
func (v *pointsItemTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *pointsItemTable) Columns() []string {
	return []string{"id", "user_id", "create_time", "log_type", "reason", "amount", "balance", "order_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *pointsItemTable) NewStruct() reform.Struct {
	return new(PointsItem)
}

// NewRecord makes a new record for that table.
func (v *pointsItemTable) NewRecord() reform.Record {
	return new(PointsItem)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *pointsItemTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// PointsItemTable represents cc_member_credit_log view or table in SQL database.
var PointsItemTable = &pointsItemTable{
	s: parse.StructInfo{Type: "PointsItem", SQLSchema: "", SQLName: "cc_member_credit_log", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "UserID", PKType: "", Column: "user_id"}, {Name: "CreatedAt", PKType: "", Column: "create_time"}, {Name: "Type", PKType: "", Column: "log_type"}, {Name: "Reason", PKType: "", Column: "reason"}, {Name: "Amount", PKType: "", Column: "amount"}, {Name: "Balance", PKType: "", Column: "balance"}, {Name: "OrderID", PKType: "", Column: "order_id"}}, PKFieldIndex: 0},
	z: new(PointsItem).Values(),
}

// String returns a string representation of this struct or record.
func (s PointsItem) String() string {
	res := make([]string, 8)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[3] = "Type: " + reform.Inspect(s.Type, true)
	res[4] = "Reason: " + reform.Inspect(s.Reason, true)
	res[5] = "Amount: " + reform.Inspect(s.Amount, true)
	res[6] = "Balance: " + reform.Inspect(s.Balance, true)
	res[7] = "OrderID: " + reform.Inspect(s.OrderID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *PointsItem) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.CreatedAt,
		s.Type,
		s.Reason,
		s.Amount,
		s.Balance,
		s.OrderID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *PointsItem) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.CreatedAt,
		&s.Type,
		&s.Reason,
		&s.Amount,
		&s.Balance,
		&s.OrderID,
	}
}

// View returns View object for that struct.
func (s *PointsItem) View() reform.View {
	return PointsItemTable
}

// Table returns Table object for that record.
func (s *PointsItem) Table() reform.Table {
	return PointsItemTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *PointsItem) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *PointsItem) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *PointsItem) HasPK() bool {
	return s.ID != PointsItemTable.z[PointsItemTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *PointsItem) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = PointsItemTable
	_ reform.Struct = new(PointsItem)
	_ reform.Table  = PointsItemTable
	_ reform.Record = new(PointsItem)
	_ fmt.Stringer  = new(PointsItem)
)

func init() {
	parse.AssertUpToDate(&CapitalFlowTable.s, new(CapitalFlow))
	CapitalFlowTable.ViewBase = reform.NewViewBase(&CapitalFlowTable.s)
	parse.AssertUpToDate(&PointsItemTable.s, new(PointsItem))
	PointsItemTable.ViewBase = reform.NewViewBase(&PointsItemTable.s)
}
