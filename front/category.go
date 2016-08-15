//go:generate reform
package front

//reform:cc_product_category
type Category struct {
	ID       uint   `reform:"cate_id,pk"`
	ParentID uint   `reform:"parent_id"`
	Name     string `reform:"cate_name"`
	Img      string `reform:"img"`
	Pos      int64  `reform:"ordering"`
}
