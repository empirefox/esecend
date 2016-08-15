package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type specialTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *specialTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_special").
func (v *specialTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *specialTable) Columns() []string {
	return []string{"special_id", "special_name", "pos"}
}

// NewStruct makes a new struct for that view or table.
func (v *specialTable) NewStruct() reform.Struct {
	return new(Special)
}

// NewRecord makes a new record for that table.
func (v *specialTable) NewRecord() reform.Record {
	return new(Special)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *specialTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// SpecialTable represents cc_special view or table in SQL database.
var SpecialTable = &specialTable{
	s: parse.StructInfo{Type: "Special", SQLSchema: "", SQLName: "cc_special", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "special_id"}, {Name: "Name", PKType: "", Column: "special_name"}, {Name: "Pos", PKType: "", Column: "pos"}}, PKFieldIndex: 0},
	z: new(Special).Values(),
}

// String returns a string representation of this struct or record.
func (s Special) String() string {
	res := make([]string, 3)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	res[2] = "Pos: " + reform.Inspect(s.Pos, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Special) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
		s.Pos,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Special) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
		&s.Pos,
	}
}

// View returns View object for that struct.
func (s *Special) View() reform.View {
	return SpecialTable
}

// Table returns Table object for that record.
func (s *Special) Table() reform.Table {
	return SpecialTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Special) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Special) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Special) HasPK() bool {
	return s.ID != SpecialTable.z[SpecialTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Special) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = SpecialTable
	_ reform.Struct = new(Special)
	_ reform.Table  = SpecialTable
	_ reform.Record = new(Special)
	_ fmt.Stringer  = new(Special)
)

type productSpecialTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *productSpecialTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product_special").
func (v *productSpecialTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *productSpecialTable) Columns() []string {
	return []string{"id", "special_id", "product_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *productSpecialTable) NewStruct() reform.Struct {
	return new(ProductSpecial)
}

// NewRecord makes a new record for that table.
func (v *productSpecialTable) NewRecord() reform.Record {
	return new(ProductSpecial)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *productSpecialTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ProductSpecialTable represents cc_product_special view or table in SQL database.
var ProductSpecialTable = &productSpecialTable{
	s: parse.StructInfo{Type: "ProductSpecial", SQLSchema: "", SQLName: "cc_product_special", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "SpecialID", PKType: "", Column: "special_id"}, {Name: "ProductID", PKType: "", Column: "product_id"}}, PKFieldIndex: 0},
	z: new(ProductSpecial).Values(),
}

// String returns a string representation of this struct or record.
func (s ProductSpecial) String() string {
	res := make([]string, 3)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "SpecialID: " + reform.Inspect(s.SpecialID, true)
	res[2] = "ProductID: " + reform.Inspect(s.ProductID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *ProductSpecial) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.SpecialID,
		s.ProductID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *ProductSpecial) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.SpecialID,
		&s.ProductID,
	}
}

// View returns View object for that struct.
func (s *ProductSpecial) View() reform.View {
	return ProductSpecialTable
}

// Table returns Table object for that record.
func (s *ProductSpecial) Table() reform.Table {
	return ProductSpecialTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductSpecial) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductSpecial) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *ProductSpecial) HasPK() bool {
	return s.ID != ProductSpecialTable.z[ProductSpecialTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *ProductSpecial) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = ProductSpecialTable
	_ reform.Struct = new(ProductSpecial)
	_ reform.Table  = ProductSpecialTable
	_ reform.Record = new(ProductSpecial)
	_ fmt.Stringer  = new(ProductSpecial)
)

func init() {
	parse.AssertUpToDate(&SpecialTable.s, new(Special))
	SpecialTable.ViewBase = reform.NewViewBase(&SpecialTable.s)
	parse.AssertUpToDate(&ProductSpecialTable.s, new(ProductSpecial))
	ProductSpecialTable.ViewBase = reform.NewViewBase(&ProductSpecialTable.s)
}
