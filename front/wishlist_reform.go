package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type wishItemTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *wishItemTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_wishlist").
func (v *wishItemTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *wishItemTable) Columns() []string {
	return []string{"id", "user_id", "created_at", "name", "img", "price", "product_id"}
}

// NewStruct makes a new struct for that view or table.
func (v *wishItemTable) NewStruct() reform.Struct {
	return new(WishItem)
}

// NewRecord makes a new record for that table.
func (v *wishItemTable) NewRecord() reform.Record {
	return new(WishItem)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *wishItemTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// WishItemTable represents cc_wishlist view or table in SQL database.
var WishItemTable = &wishItemTable{
	s: parse.StructInfo{Type: "WishItem", SQLSchema: "", SQLName: "cc_wishlist", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "UserID", PKType: "", Column: "user_id"}, {Name: "CreatedAt", PKType: "", Column: "created_at"}, {Name: "Name", PKType: "", Column: "name"}, {Name: "Img", PKType: "", Column: "img"}, {Name: "Price", PKType: "", Column: "price"}, {Name: "ProductID", PKType: "", Column: "product_id"}}, PKFieldIndex: 0},
	z: new(WishItem).Values(),
}

// String returns a string representation of this struct or record.
func (s WishItem) String() string {
	res := make([]string, 7)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[3] = "Name: " + reform.Inspect(s.Name, true)
	res[4] = "Img: " + reform.Inspect(s.Img, true)
	res[5] = "Price: " + reform.Inspect(s.Price, true)
	res[6] = "ProductID: " + reform.Inspect(s.ProductID, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *WishItem) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.CreatedAt,
		s.Name,
		s.Img,
		s.Price,
		s.ProductID,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *WishItem) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.CreatedAt,
		&s.Name,
		&s.Img,
		&s.Price,
		&s.ProductID,
	}
}

// View returns View object for that struct.
func (s *WishItem) View() reform.View {
	return WishItemTable
}

// Table returns Table object for that record.
func (s *WishItem) Table() reform.Table {
	return WishItemTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *WishItem) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *WishItem) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *WishItem) HasPK() bool {
	return s.ID != WishItemTable.z[WishItemTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *WishItem) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = WishItemTable
	_ reform.Struct = new(WishItem)
	_ reform.Table  = WishItemTable
	_ reform.Record = new(WishItem)
	_ fmt.Stringer  = new(WishItem)
)

func init() {
	parse.AssertUpToDate(&WishItemTable.s, new(WishItem))
	WishItemTable.ViewBase = reform.NewViewBase(&WishItemTable.s)
}
