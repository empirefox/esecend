//go:generate reform
package front

import "github.com/empirefox/reform"

//reform:cc_product_attribute
type ProductAttr struct {
	ID      uint   `reform:"id,pk"`
	Value   string `reform:"attribute_value"`
	GroupID uint   `reform:"att_group_id"`
	Pos     int64  `reform:"pos"`
}

//reform:cc_product_attribute_group
type ProductAttrGroup struct {
	ID   uint   `reform:"attribute_cate_id,pk"`
	Name string `reform:"attribute_cate_name"`
	Pos  int64  `reform:"pos"`
}

type ProductAttrsResponse struct {
	Groups   []reform.Struct // ProductAttrGroup
	Attrs    []reform.Struct // ProductAttr
	Specials []reform.Struct // Special
}
