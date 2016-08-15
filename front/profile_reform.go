package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type profileView struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *profileView) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_profile").
func (v *profileView) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *profileView) Columns() []string {
	return []string{"phone", "free_delivery_line", "default_head_image"}
}

// NewStruct makes a new struct for that view or table.
func (v *profileView) NewStruct() reform.Struct {
	return new(Profile)
}

// ProfileView represents cc_profile view or table in SQL database.
var ProfileView = &profileView{
	s: parse.StructInfo{Type: "Profile", SQLSchema: "", SQLName: "cc_profile", Fields: []parse.FieldInfo{{Name: "Phone", PKType: "", Column: "phone"}, {Name: "FreeDeliverLine", PKType: "", Column: "free_delivery_line"}, {Name: "DefaultHeadImage", PKType: "", Column: "default_head_image"}}, PKFieldIndex: -1},
	z: new(Profile).Values(),
}

// String returns a string representation of this struct or record.
func (s Profile) String() string {
	res := make([]string, 3)
	res[0] = "Phone: " + reform.Inspect(s.Phone, true)
	res[1] = "FreeDeliverLine: " + reform.Inspect(s.FreeDeliverLine, true)
	res[2] = "DefaultHeadImage: " + reform.Inspect(s.DefaultHeadImage, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Profile) Values() []interface{} {
	return []interface{}{
		s.Phone,
		s.FreeDeliverLine,
		s.DefaultHeadImage,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Profile) Pointers() []interface{} {
	return []interface{}{
		&s.Phone,
		&s.FreeDeliverLine,
		&s.DefaultHeadImage,
	}
}

// View returns View object for that struct.
func (s *Profile) View() reform.View {
	return ProfileView
}

// check interfaces
var (
	_ reform.View   = ProfileView
	_ reform.Struct = new(Profile)
	_ fmt.Stringer  = new(Profile)
)

func init() {
	parse.AssertUpToDate(&ProfileView.s, new(Profile))
	ProfileView.ViewBase = reform.NewViewBase(&ProfileView.s)
}
