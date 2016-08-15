package front

// generated with github.com/empirefox/reform

import (
	"fmt"
	"strings"

	"github.com/empirefox/reform"
	"github.com/empirefox/reform/parse"
)

type evalItemView struct {
	*reform.ViewBase
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *evalItemView) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("cc_order_item").
func (v *evalItemView) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *evalItemView) Columns() []string {
	return []string{"comment_content", "comment_time", "user_name", "starts", "rate_fit", "rate_serve", "rate_deliver"}
}

// NewStruct makes a new struct for that view or table.
func (v *evalItemView) NewStruct() reform.Struct {
	return new(EvalItem)
}

// EvalItemView represents cc_order_item view or table in SQL database.
var EvalItemView = &evalItemView{
	s: parse.StructInfo{Type: "EvalItem", SQLSchema: "", SQLName: "cc_order_item", Fields: []parse.FieldInfo{{Name: "Eval", PKType: "", Column: "comment_content"}, {Name: "EvalAt", PKType: "", Column: "comment_time"}, {Name: "EvalName", PKType: "", Column: "user_name"}, {Name: "RateStar", PKType: "", Column: "starts"}, {Name: "RateFit", PKType: "", Column: "rate_fit"}, {Name: "RateServe", PKType: "", Column: "rate_serve"}, {Name: "RateDeliver", PKType: "", Column: "rate_deliver"}}, PKFieldIndex: -1},
	z: new(EvalItem).Values(),
}

// String returns a string representation of this struct or record.
func (s EvalItem) String() string {
	res := make([]string, 7)
	res[0] = "Eval: " + reform.Inspect(s.Eval, true)
	res[1] = "EvalAt: " + reform.Inspect(s.EvalAt, true)
	res[2] = "EvalName: " + reform.Inspect(s.EvalName, true)
	res[3] = "RateStar: " + reform.Inspect(s.RateStar, true)
	res[4] = "RateFit: " + reform.Inspect(s.RateFit, true)
	res[5] = "RateServe: " + reform.Inspect(s.RateServe, true)
	res[6] = "RateDeliver: " + reform.Inspect(s.RateDeliver, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *EvalItem) Values() []interface{} {
	return []interface{}{
		s.Eval,
		s.EvalAt,
		s.EvalName,
		s.RateStar,
		s.RateFit,
		s.RateServe,
		s.RateDeliver,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *EvalItem) Pointers() []interface{} {
	return []interface{}{
		&s.Eval,
		&s.EvalAt,
		&s.EvalName,
		&s.RateStar,
		&s.RateFit,
		&s.RateServe,
		&s.RateDeliver,
	}
}

// View returns View object for that struct.
func (s *EvalItem) View() reform.View {
	return EvalItemView
}

// check interfaces
var (
	_ reform.View   = EvalItemView
	_ reform.Struct = new(EvalItem)
	_ fmt.Stringer  = new(EvalItem)
)

func init() {
	parse.AssertUpToDate(&EvalItemView.s, new(EvalItem))
	EvalItemView.ViewBase = reform.NewViewBase(&EvalItemView.s)
}
