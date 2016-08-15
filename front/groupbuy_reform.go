package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type groupBuyItemTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *groupBuyItemTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_group_buy").
func (v *groupBuyItemTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *groupBuyItemTable) Columns() []string {
	return []string{"id", "img", "title", "reason", "price", "start", "end", "sku_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *groupBuyItemTable) NewStruct() reform.Struct {
	return new(GroupBuyItem)
}

// NewRecord makes a new record for that table.
func (v *groupBuyItemTable) NewRecord() reform.Record {
	return new(GroupBuyItem)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *groupBuyItemTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// GroupBuyItemTable represents cc_group_buy view or table in SQL database.
var GroupBuyItemTable = &groupBuyItemTable{
	s: parse.StructInfo{Type: "GroupBuyItem", SQLSchema: "", SQLName: "cc_group_buy", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Title", PKType: "", Column: "title"}, {Name: "Reason", PKType: "", Column: "reason"}, {Name: "Price", PKType: "", Column: "price"}, {Name: "Start", PKType: "", Column: "start"}, {Name: "End", PKType: "", Column: "end"}, {Name: "SkuID", PKType: "", Column: "sku_id"}}, PKFieldIndex: 0},
	z: new(GroupBuyItem).Values(),
}

// String returns a string representation of this struct or record.
func (s GroupBuyItem) String() string {
	res := make([]string, 8)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Img: " + reform.Inspect(s.Img, true)
	res[2] = "Title: " + reform.Inspect(s.Title, true)
	res[3] = "Reason: " + reform.Inspect(s.Reason, true)
	res[4] = "Price: " + reform.Inspect(s.Price, true)
	res[5] = "Start: " + reform.Inspect(s.Start, true)
	res[6] = "End: " + reform.Inspect(s.End, true)
	res[7] = "SkuID: " + reform.Inspect(s.SkuID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *GroupBuyItem) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Img,
		s.Title,
		s.Reason,
		s.Price,
		s.Start,
		s.End,
		s.SkuID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *GroupBuyItem) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Img,
		&s.Title,
		&s.Reason,
		&s.Price,
		&s.Start,
		&s.End,
		&s.SkuID,
	}
}

// View returns View object for that struct.
func (s *GroupBuyItem) View() reform.View {
	return GroupBuyItemTable
}

// Table returns Table object for that record.
func (s *GroupBuyItem) Table() reform.Table {
	return GroupBuyItemTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *GroupBuyItem) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *GroupBuyItem) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *GroupBuyItem) HasPK() bool {
	return s.ID != GroupBuyItemTable.z[GroupBuyItemTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *GroupBuyItem) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = GroupBuyItemTable
	_ reform.Struct = new(GroupBuyItem)
	_ reform.Table  = GroupBuyItemTable
	_ reform.Record = new(GroupBuyItem)
	_ fmt.Stringer  = new(GroupBuyItem)
)

func init() {
	parse.AssertUpToDate(&GroupBuyItemTable.s, new(GroupBuyItem))
	GroupBuyItemTable.ViewBase = reform.NewViewBase(&GroupBuyItemTable.s)
}
