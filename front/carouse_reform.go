package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type carouselItemTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *carouselItemTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_carousel").
func (v *carouselItemTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *carouselItemTable) Columns() []string {
	return []string{"id", "img", "link", "product_id", "category_id", "pos"}
}

// NewStruct makes a new struct for that view or table.
func (v *carouselItemTable) NewStruct() reform.Struct {
	return new(CarouselItem)
}

// NewRecord makes a new record for that table.
func (v *carouselItemTable) NewRecord() reform.Record {
	return new(CarouselItem)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *carouselItemTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// CarouselItemTable represents cc_carousel view or table in SQL database.
var CarouselItemTable = &carouselItemTable{
	s: parse.StructInfo{Type: "CarouselItem", SQLSchema: "", SQLName: "cc_carousel", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Link", PKType: "", Column: "link"}, {Name: "ProductID", PKType: "", Column: "product_id"}, {Name: "CategoryID", PKType: "", Column: "category_id"}, {Name: "Pos", PKType: "", Column: "pos"}}, PKFieldIndex: 0},
	z: new(CarouselItem).Values(),
}

// String returns a string representation of this struct or record.
func (s CarouselItem) String() string {
	res := make([]string, 6)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Img: " + reform.Inspect(s.Img, true)
	res[2] = "Link: " + reform.Inspect(s.Link, true)
	res[3] = "ProductID: " + reform.Inspect(s.ProductID, true)
	res[4] = "CategoryID: " + reform.Inspect(s.CategoryID, true)
	res[5] = "Pos: " + reform.Inspect(s.Pos, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *CarouselItem) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Img,
		s.Link,
		s.ProductID,
		s.CategoryID,
		s.Pos,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *CarouselItem) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Img,
		&s.Link,
		&s.ProductID,
		&s.CategoryID,
		&s.Pos,
	}
}

// View returns View object for that struct.
func (s *CarouselItem) View() reform.View {
	return CarouselItemTable
}

// Table returns Table object for that record.
func (s *CarouselItem) Table() reform.Table {
	return CarouselItemTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *CarouselItem) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *CarouselItem) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *CarouselItem) HasPK() bool {
	return s.ID != CarouselItemTable.z[CarouselItemTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *CarouselItem) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = CarouselItemTable
	_ reform.Struct = new(CarouselItem)
	_ reform.Table  = CarouselItemTable
	_ reform.Record = new(CarouselItem)
	_ fmt.Stringer  = new(CarouselItem)
)

func init() {
	parse.AssertUpToDate(&CarouselItemTable.s, new(CarouselItem))
	CarouselItemTable.ViewBase = reform.NewViewBase(&CarouselItemTable.s)
}
