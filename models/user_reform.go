package models

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type userTable struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *userTable) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_member").
func (v *userTable) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *userTable) Columns() []string {
	return []string{"id", "open_id", "level_id", "privilege", "phone", "parent_id", "create_date", "update_date", "last_login", "name", "sex", "city", "province", "avatar", "union_id", "refresh_token", "paykey"}
}

// NewStruct makes a new struct for that view or table.
func (v *userTable) NewStruct() reform.Struct {
	return new(User)
}

// NewRecord makes a new record for that table.
func (v *userTable) NewRecord() reform.Record {
	return new(User)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *userTable) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// UserTable represents cc_member view or table in SQL database.
var UserTable = &userTable{
	s: parse.StructInfo{Type: "User", SQLSchema: "", SQLName: "cc_member", Fields: []parse.FieldInfo{{Name: "ID", PKType: "uint", Column: "id"}, {Name: "OpenId", PKType: "", Column: "open_id"}, {Name: "LevelID", PKType: "", Column: "level_id"}, {Name: "Privilege", PKType: "", Column: "privilege"}, {Name: "Phone", PKType: "", Column: "phone"}, {Name: "Recommended", PKType: "", Column: "parent_id"}, {Name: "CreatedAt", PKType: "", Column: "create_date"}, {Name: "UpdatedAt", PKType: "", Column: "update_date"}, {Name: "SigninAt", PKType: "", Column: "last_login"}, {Name: "Nickname", PKType: "", Column: "name"}, {Name: "Sex", PKType: "", Column: "sex"}, {Name: "City", PKType: "", Column: "city"}, {Name: "Province", PKType: "", Column: "province"}, {Name: "HeadImageURL", PKType: "", Column: "avatar"}, {Name: "UnionId", PKType: "", Column: "union_id"}, {Name: "RefreshToken", PKType: "", Column: "refresh_token"}, {Name: "Paykey", PKType: "", Column: "paykey"}}, PKFieldIndex: 0},
	z: new(User).Values(),
}

// String returns a string representation of this struct or record.
func (s User) String() string {
	res := make([]string, 17)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "OpenId: " + reform.Inspect(s.OpenId, true)
	res[2] = "LevelID: " + reform.Inspect(s.LevelID, true)
	res[3] = "Privilege: " + reform.Inspect(s.Privilege, true)
	res[4] = "Phone: " + reform.Inspect(s.Phone, true)
	res[5] = "Recommended: " + reform.Inspect(s.Recommended, true)
	res[6] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[7] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	res[8] = "SigninAt: " + reform.Inspect(s.SigninAt, true)
	res[9] = "Nickname: " + reform.Inspect(s.Nickname, true)
	res[10] = "Sex: " + reform.Inspect(s.Sex, true)
	res[11] = "City: " + reform.Inspect(s.City, true)
	res[12] = "Province: " + reform.Inspect(s.Province, true)
	res[13] = "HeadImageURL: " + reform.Inspect(s.HeadImageURL, true)
	res[14] = "UnionId: " + reform.Inspect(s.UnionId, true)
	res[15] = "RefreshToken: " + reform.Inspect(s.RefreshToken, true)
	res[16] = "Paykey: " + reform.Inspect(s.Paykey, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *User) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.OpenId,
		s.LevelID,
		s.Privilege,
		s.Phone,
		s.Recommended,
		s.CreatedAt,
		s.UpdatedAt,
		s.SigninAt,
		s.Nickname,
		s.Sex,
		s.City,
		s.Province,
		s.HeadImageURL,
		s.UnionId,
		s.RefreshToken,
		s.Paykey,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *User) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.OpenId,
		&s.LevelID,
		&s.Privilege,
		&s.Phone,
		&s.Recommended,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.SigninAt,
		&s.Nickname,
		&s.Sex,
		&s.City,
		&s.Province,
		&s.HeadImageURL,
		&s.UnionId,
		&s.RefreshToken,
		&s.Paykey,
	}
}

// View returns View object for that struct.
func (s *User) View() reform.View {
	return UserTable
}

// Table returns Table object for that record.
func (s *User) Table() reform.Table {
	return UserTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *User) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *User) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *User) HasPK() bool {
	return s.ID != UserTable.z[UserTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *User) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint(i64)
	} else {
		s.ID = pk.(uint)
	}
}

// check interfaces
var (
	_ reform.View   = UserTable
	_ reform.Struct = new(User)
	_ reform.Table  = UserTable
	_ reform.Record = new(User)
	_ fmt.Stringer  = new(User)
)

func init() {
	parse.AssertUpToDate(&UserTable.s, new(User))
	UserTable.ViewBase = reform.NewViewBase(&UserTable.s)
}
