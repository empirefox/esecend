package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type productTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *productTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product").
func (v *productTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *productTable) Columns() []string {
	return []string{"product_id", "product_name", "img", "intro", "detail", "saleCount", "create_date", "time_sale", "time_shelfoff", "cate_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *productTable) NewStruct() reform.Struct {
	return new(Product)
}

// NewRecord makes a new record for that table.
func (v *productTable) NewRecord() reform.Record {
	return new(Product)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *productTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ProductTable represents cc_product view or table in SQL database.
var ProductTable = &productTable{
	s: parse.StructInfo{Type: "Product", SQLSchema: "", SQLName: "cc_product", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "product_id"}, {Name: "Name", PKType: "", Column: "product_name"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Intro", PKType: "", Column: "intro"}, {Name: "Detail", PKType: "", Column: "detail"}, {Name: "Saled", PKType: "", Column: "saleCount"}, {Name: "CreatedAt", PKType: "", Column: "create_date"}, {Name: "SaledAt", PKType: "", Column: "time_sale"}, {Name: "ShelfOffAt", PKType: "", Column: "time_shelfoff"}, {Name: "CategoryID", PKType: "", Column: "cate_id"}}, PKFieldIndex: 0},
	z: new(Product).Values(),
}

// String returns a string representation of this struct or record.
func (s Product) String() string {
	res := make([]string, 10)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	res[2] = "Img: " + reform.Inspect(s.Img, true)
	res[3] = "Intro: " + reform.Inspect(s.Intro, true)
	res[4] = "Detail: " + reform.Inspect(s.Detail, true)
	res[5] = "Saled: " + reform.Inspect(s.Saled, true)
	res[6] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[7] = "SaledAt: " + reform.Inspect(s.SaledAt, true)
	res[8] = "ShelfOffAt: " + reform.Inspect(s.ShelfOffAt, true)
	res[9] = "CategoryID: " + reform.Inspect(s.CategoryID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Product) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
		s.Img,
		s.Intro,
		s.Detail,
		s.Saled,
		s.CreatedAt,
		s.SaledAt,
		s.ShelfOffAt,
		s.CategoryID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Product) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
		&s.Img,
		&s.Intro,
		&s.Detail,
		&s.Saled,
		&s.CreatedAt,
		&s.SaledAt,
		&s.ShelfOffAt,
		&s.CategoryID,
	}
}

// View returns View object for that struct.
func (s *Product) View() reform.View {
	return ProductTable
}

// Table returns Table object for that record.
func (s *Product) Table() reform.Table {
	return ProductTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Product) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Product) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Product) HasPK() bool {
	return s.ID != ProductTable.z[ProductTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Product) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = ProductTable
	_ reform.Struct = new(Product)
	_ reform.Table  = ProductTable
	_ reform.Record = new(Product)
	_ fmt.Stringer  = new(Product)
)

type productAttrIdTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *productAttrIdTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product_sku_att").
func (v *productAttrIdTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *productAttrIdTable) Columns() []string {
	return []string{"id", "sku_id", "att_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *productAttrIdTable) NewStruct() reform.Struct {
	return new(ProductAttrId)
}

// NewRecord makes a new record for that table.
func (v *productAttrIdTable) NewRecord() reform.Record {
	return new(ProductAttrId)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *productAttrIdTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ProductAttrIdTable represents cc_product_sku_att view or table in SQL database.
var ProductAttrIdTable = &productAttrIdTable{
	s: parse.StructInfo{Type: "ProductAttrId", SQLSchema: "", SQLName: "cc_product_sku_att", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "SkuID", PKType: "", Column: "sku_id"}, {Name: "AttrID", PKType: "", Column: "att_id"}}, PKFieldIndex: 0},
	z: new(ProductAttrId).Values(),
}

// String returns a string representation of this struct or record.
func (s ProductAttrId) String() string {
	res := make([]string, 3)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "SkuID: " + reform.Inspect(s.SkuID, true)
	res[2] = "AttrID: " + reform.Inspect(s.AttrID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *ProductAttrId) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.SkuID,
		s.AttrID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *ProductAttrId) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.SkuID,
		&s.AttrID,
	}
}

// View returns View object for that struct.
func (s *ProductAttrId) View() reform.View {
	return ProductAttrIdTable
}

// Table returns Table object for that record.
func (s *ProductAttrId) Table() reform.Table {
	return ProductAttrIdTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductAttrId) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *ProductAttrId) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *ProductAttrId) HasPK() bool {
	return s.ID != ProductAttrIdTable.z[ProductAttrIdTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *ProductAttrId) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = ProductAttrIdTable
	_ reform.Struct = new(ProductAttrId)
	_ reform.Table  = ProductAttrIdTable
	_ reform.Record = new(ProductAttrId)
	_ fmt.Stringer  = new(ProductAttrId)
)

func init() {
	parse.AssertUpToDate(&ProductTable.s, new(Product))
	ProductTable.ViewBase = reform.NewViewBase(&ProductTable.s)
	parse.AssertUpToDate(&ProductAttrIdTable.s, new(ProductAttrId))
	ProductAttrIdTable.ViewBase = reform.NewViewBase(&ProductAttrIdTable.s)
}
