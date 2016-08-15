package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type skuTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *skuTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_product_sku").
func (v *skuTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *skuTable) Columns() []string {
	return []string{"sku_id", "stock", "img", "sale_price", "market_price", "freight", "product_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *skuTable) NewStruct() reform.Struct {
	return new(Sku)
}

// NewRecord makes a new record for that table.
func (v *skuTable) NewRecord() reform.Record {
	return new(Sku)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *skuTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// SkuTable represents cc_product_sku view or table in SQL database.
var SkuTable = &skuTable{
	s: parse.StructInfo{Type: "Sku", SQLSchema: "", SQLName: "cc_product_sku", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "sku_id"}, {Name: "Stock", PKType: "", Column: "stock"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "SalePrice", PKType: "", Column: "sale_price"}, {Name: "MarketPrice", PKType: "", Column: "market_price"}, {Name: "Freight", PKType: "", Column: "freight"}, {Name: "ProductID", PKType: "", Column: "product_id"}}, PKFieldIndex: 0},
	z: new(Sku).Values(),
}

// String returns a string representation of this struct or record.
func (s Sku) String() string {
	res := make([]string, 7)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Stock: " + reform.Inspect(s.Stock, true)
	res[2] = "Img: " + reform.Inspect(s.Img, true)
	res[3] = "SalePrice: " + reform.Inspect(s.SalePrice, true)
	res[4] = "MarketPrice: " + reform.Inspect(s.MarketPrice, true)
	res[5] = "Freight: " + reform.Inspect(s.Freight, true)
	res[6] = "ProductID: " + reform.Inspect(s.ProductID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Sku) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Stock,
		s.Img,
		s.SalePrice,
		s.MarketPrice,
		s.Freight,
		s.ProductID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Sku) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Stock,
		&s.Img,
		&s.SalePrice,
		&s.MarketPrice,
		&s.Freight,
		&s.ProductID,
	}
}

// View returns View object for that struct.
func (s *Sku) View() reform.View {
	return SkuTable
}

// Table returns Table object for that record.
func (s *Sku) Table() reform.Table {
	return SkuTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Sku) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Sku) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Sku) HasPK() bool {
	return s.ID != SkuTable.z[SkuTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Sku) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = SkuTable
	_ reform.Struct = new(Sku)
	_ reform.Table  = SkuTable
	_ reform.Record = new(Sku)
	_ fmt.Stringer  = new(Sku)
)

func init() {
	parse.AssertUpToDate(&SkuTable.s, new(Sku))
	SkuTable.ViewBase = reform.NewViewBase(&SkuTable.s)
}
