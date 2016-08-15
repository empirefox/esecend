package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type categoryTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *categoryTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product_category").
func (v *categoryTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *categoryTable) Columns() []string {
	return []string{"cate_id", "parent_id", "cate_name", "img", "ordering"}
}

// NewStruct makes a new struct for that view or table.
func (v *categoryTable) NewStruct() reform.Struct {
	return new(Category)
}

// NewRecord makes a new record for that table.
func (v *categoryTable) NewRecord() reform.Record {
	return new(Category)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *categoryTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// CategoryTable represents cc_product_category view or table in SQL database.
var CategoryTable = &categoryTable{
	s: parse.StructInfo{Type: "Category", SQLSchema: "", SQLName: "cc_product_category", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "cate_id"}, {Name: "ParentID", PKType: "", Column: "parent_id"}, {Name: "Name", PKType: "", Column: "cate_name"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Pos", PKType: "", Column: "ordering"}}, PKFieldIndex: 0},
	z: new(Category).Values(),
}

// String returns a string representation of this struct or record.
func (s Category) String() string {
	res := make([]string, 5)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "ParentID: " + reform.Inspect(s.ParentID, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Img: " + reform.Inspect(s.Img, true)
	res[4] = "Pos: " + reform.Inspect(s.Pos, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Category) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.ParentID,
		s.Name,
		s.Img,
		s.Pos,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Category) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.ParentID,
		&s.Name,
		&s.Img,
		&s.Pos,
	}
}

// View returns View object for that struct.
func (s *Category) View() reform.View {
	return CategoryTable
}

// Table returns Table object for that record.
func (s *Category) Table() reform.Table {
	return CategoryTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Category) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Category) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Category) HasPK() bool {
	return s.ID != CategoryTable.z[CategoryTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Category) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = CategoryTable
	_ reform.Struct = new(Category)
	_ reform.Table  = CategoryTable
	_ reform.Record = new(Category)
	_ fmt.Stringer  = new(Category)
)

func init() {
	parse.AssertUpToDate(&CategoryTable.s, new(Category))
	CategoryTable.ViewBase = reform.NewViewBase(&CategoryTable.s)
}
