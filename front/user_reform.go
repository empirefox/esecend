package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type userLevelTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *userLevelTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_member_level").
func (v *userLevelTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *userLevelTable) Columns() []string {
	return []string{"id", "name"}
}

// NewStruct makes a new struct for that view or table.
func (v *userLevelTable) NewStruct() reform.Struct {
	return new(UserLevel)
}

// NewRecord makes a new record for that table.
func (v *userLevelTable) NewRecord() reform.Record {
	return new(UserLevel)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *userLevelTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// UserLevelTable represents cc_member_level view or table in SQL database.
var UserLevelTable = &userLevelTable{
	s: parse.StructInfo{Type: "UserLevel", SQLSchema: "", SQLName: "cc_member_level", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "Name", PKType: "", Column: "name"}}, PKFieldIndex: 0},
	z: new(UserLevel).Values(),
}

// String returns a string representation of this struct or record.
func (s UserLevel) String() string {
	res := make([]string, 2)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Name: " + reform.Inspect(s.Name, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *UserLevel) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Name,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *UserLevel) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Name,
	}
}

// View returns View object for that struct.
func (s *UserLevel) View() reform.View {
	return UserLevelTable
}

// Table returns Table object for that record.
func (s *UserLevel) Table() reform.Table {
	return UserLevelTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *UserLevel) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *UserLevel) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *UserLevel) HasPK() bool {
	return s.ID != UserLevelTable.z[UserLevelTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *UserLevel) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = UserLevelTable
	_ reform.Struct = new(UserLevel)
	_ reform.Table  = UserLevelTable
	_ reform.Record = new(UserLevel)
	_ fmt.Stringer  = new(UserLevel)
)

func init() {
	parse.AssertUpToDate(&UserLevelTable.s, new(UserLevel))
	UserLevelTable.ViewBase = reform.NewViewBase(&UserLevelTable.s)
}
