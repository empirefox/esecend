package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type productAttrTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *productAttrTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product_attribute").
func (v *productAttrTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *productAttrTable) Columns() []string {
	return []string{"id", "attribute_value", "att_group_id", "pos"}
}

// NewStruct makes a new struct for that view or table.
func (v *productAttrTable) NewStruct() reform.Struct {
	return new(ProductAttr)
}

// NewRecord makes a new record for that table.
func (v *productAttrTable) NewRecord() reform.Record {
	return new(ProductAttr)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *productAttrTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ProductAttrTable represents cc_product_attribute view or table in SQL database.
var ProductAttrTable = &productAttrTable{
	s: parse.StructInfo{Type: "ProductAttr", SQLSchema: "", SQLName: "cc_product_attribute", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "Value", PKType: "", Column: "attribute_value"}, {Name: "GroupID", PKType: "", Column: "att_group_id"}, {Name: "Pos", PKType: "", Column: "pos"}}, PKFieldIndex: 0},
	z: new(ProductAttr).Values(),
}

// String returns a string representation of this struct or record.
func (s ProductAttr) String() string {
	res := make([]string, 4)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Value: " + reform.Inspect(s.Value, true)
	res[2] = "GroupID: " + reform.Inspect(s.GroupID, true)
	res[3] = "Pos: " + reform.Inspect(s.Pos, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *ProductAttr) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Value,
		s.GroupID,
		s.Pos,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *ProductAttr) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Value,
		&s.GroupID,
		&s.Pos,
	}
}

// View returns View object for that struct.
func (s *ProductAttr) View() reform.View {
	return ProductAttrTable
}

// Table returns Table object for that record.
func (s *ProductAttr) Table() reform.Table {
	return ProductAttrTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductAttr) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductAttr) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *ProductAttr) HasPK() bool {
	return s.ID != ProductAttrTable.z[ProductAttrTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *ProductAttr) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = ProductAttrTable
	_ reform.Struct = new(ProductAttr)
	_ reform.Table  = ProductAttrTable
	_ reform.Record = new(ProductAttr)
	_ fmt.Stringer  = new(ProductAttr)
)

type productAttrGroupTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *productAttrGroupTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product_attribute_group").
func (v *productAttrGroupTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *productAttrGroupTable) Columns() []string {
	return []string{"attribute_cate_id", "attribute_cate_name", "pos"}
}

// NewStruct makes a new struct for that view or table.
func (v *productAttrGroupTable) NewStruct() reform.Struct {
	return new(ProductAttrGroup)
}

// NewRecord makes a new record for that table.
func (v *productAttrGroupTable) NewRecord() reform.Record {
	return new(ProductAttrGroup)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *productAttrGroupTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ProductAttrGroupTable represents cc_product_attribute_group view or table in SQL database.
var ProductAttrGroupTable = &productAttrGroupTable{
	s: parse.StructInfo{Type: "ProductAttrGroup", SQLSchema: "", SQLName: "cc_product_attribute_group", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "attribute_cate_id"}, {Name: "Name", PKType: "", Column: "attribute_cate_name"}, {Name: "Pos", PKType: "", Column: "pos"}}, PKFieldIndex: 0},
	z: new(ProductAttrGroup).Values(),
}

// String returns a string representation of this struct or record.
func (s ProductAttrGroup) String() string {
	res := make([]string, 3)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	res[2] = "Pos: " + reform.Inspect(s.Pos, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *ProductAttrGroup) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
		s.Pos,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *ProductAttrGroup) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
		&s.Pos,
	}
}

// View returns View object for that struct.
func (s *ProductAttrGroup) View() reform.View {
	return ProductAttrGroupTable
}

// Table returns Table object for that record.
func (s *ProductAttrGroup) Table() reform.Table {
	return ProductAttrGroupTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductAttrGroup) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductAttrGroup) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *ProductAttrGroup) HasPK() bool {
	return s.ID != ProductAttrGroupTable.z[ProductAttrGroupTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *ProductAttrGroup) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = ProductAttrGroupTable
	_ reform.Struct = new(ProductAttrGroup)
	_ reform.Table  = ProductAttrGroupTable
	_ reform.Record = new(ProductAttrGroup)
	_ fmt.Stringer  = new(ProductAttrGroup)
)

func init() {
	parse.AssertUpToDate(&ProductAttrTable.s, new(ProductAttr))
	ProductAttrTable.ViewBase = reform.NewViewBase(&ProductAttrTable.s)
	parse.AssertUpToDate(&ProductAttrGroupTable.s, new(ProductAttrGroup))
	ProductAttrGroupTable.ViewBase = reform.NewViewBase(&ProductAttrGroupTable.s)
}
