package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type addressTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *addressTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_member_address").
func (v *addressTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *addressTable) Columns() []string {
	return []string{"id", "uid", "contactor", "tel_num", "province", "city", "district", "address", "pos"}
}

// NewStruct makes a new struct for that view or table.
func (v *addressTable) NewStruct() reform.Struct {
	return new(Address)
}

// NewRecord makes a new record for that table.
func (v *addressTable) NewRecord() reform.Record {
	return new(Address)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *addressTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// AddressTable represents cc_member_address view or table in SQL database.
var AddressTable = &addressTable{
	s: parse.StructInfo{Type: "Address", SQLSchema: "", SQLName: "cc_member_address", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "UserID", PKType: "", Column: "uid"}, {Name: "Contact", PKType: "", Column: "contactor"}, {Name: "Phone", PKType: "", Column: "tel_num"}, {Name: "Province", PKType: "", Column: "province"}, {Name: "City", PKType: "", Column: "city"}, {Name: "District", PKType: "", Column: "district"}, {Name: "House", PKType: "", Column: "address"}, {Name: "Pos", PKType: "", Column: "pos"}}, PKFieldIndex: 0},
	z: new(Address).Values(),
}

// String returns a string representation of this struct or record.
func (s Address) String() string {
	res := make([]string, 9)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "Contact: " + reform.Inspect(s.Contact, true)
	res[3] = "Phone: " + reform.Inspect(s.Phone, true)
	res[4] = "Province: " + reform.Inspect(s.Province, true)
	res[5] = "City: " + reform.Inspect(s.City, true)
	res[6] = "District: " + reform.Inspect(s.District, true)
	res[7] = "House: " + reform.Inspect(s.House, true)
	res[8] = "Pos: " + reform.Inspect(s.Pos, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Address) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.Contact,
		s.Phone,
		s.Province,
		s.City,
		s.District,
		s.House,
		s.Pos,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Address) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.Contact,
		&s.Phone,
		&s.Province,
		&s.City,
		&s.District,
		&s.House,
		&s.Pos,
	}
}

// View returns View object for that struct.
func (s *Address) View() reform.View {
	return AddressTable
}

// Table returns Table object for that record.
func (s *Address) Table() reform.Table {
	return AddressTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Address) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Address) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Address) HasPK() bool {
	return s.ID != AddressTable.z[AddressTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Address) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = AddressTable
	_ reform.Struct = new(Address)
	_ reform.Table  = AddressTable
	_ reform.Record = new(Address)
	_ fmt.Stringer  = new(Address)
)

func init() {
	parse.AssertUpToDate(&AddressTable.s, new(Address))
	AddressTable.ViewBase = reform.NewViewBase(&AddressTable.s)
}
