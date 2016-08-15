package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type cartItemTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *cartItemTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_cart").
func (v *cartItemTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *cartItemTable) Columns() []string {
	return []string{"id", "user_id", "name", "img", "type", "price", "quantity", "created_at", "sku_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *cartItemTable) NewStruct() reform.Struct {
	return new(CartItem)
}

// NewRecord makes a new record for that table.
func (v *cartItemTable) NewRecord() reform.Record {
	return new(CartItem)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *cartItemTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// CartItemTable represents cc_cart view or table in SQL database.
var CartItemTable = &cartItemTable{
	s: parse.StructInfo{Type: "CartItem", SQLSchema: "", SQLName: "cc_cart", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "UserID", PKType: "", Column: "user_id"}, {Name: "Name", PKType: "", Column: "name"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Type", PKType: "", Column: "type"}, {Name: "Price", PKType: "", Column: "price"}, {Name: "Quantity", PKType: "", Column: "quantity"}, {Name: "CreatedAt", PKType: "", Column: "created_at"}, {Name: "SkuID", PKType: "", Column: "sku_id"}}, PKFieldIndex: 0},
	z: new(CartItem).Values(),
}

// String returns a string representation of this struct or record.
func (s CartItem) String() string {
	res := make([]string, 9)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Img: " + reform.Inspect(s.Img, true)
	res[4] = "Type: " + reform.Inspect(s.Type, true)
	res[5] = "Price: " + reform.Inspect(s.Price, true)
	res[6] = "Quantity: " + reform.Inspect(s.Quantity, true)
	res[7] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[8] = "SkuID: " + reform.Inspect(s.SkuID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *CartItem) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.Name,
		s.Img,
		s.Type,
		s.Price,
		s.Quantity,
		s.CreatedAt,
		s.SkuID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *CartItem) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.Name,
		&s.Img,
		&s.Type,
		&s.Price,
		&s.Quantity,
		&s.CreatedAt,
		&s.SkuID,
	}
}

// View returns View object for that struct.
func (s *CartItem) View() reform.View {
	return CartItemTable
}

// Table returns Table object for that record.
func (s *CartItem) Table() reform.Table {
	return CartItemTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *CartItem) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *CartItem) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *CartItem) HasPK() bool {
	return s.ID != CartItemTable.z[CartItemTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *CartItem) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = CartItemTable
	_ reform.Struct = new(CartItem)
	_ reform.Table  = CartItemTable
	_ reform.Record = new(CartItem)
	_ fmt.Stringer  = new(CartItem)
)

func init() {
	parse.AssertUpToDate(&CartItemTable.s, new(CartItem))
	CartItemTable.ViewBase = reform.NewViewBase(&CartItemTable.s)
}
