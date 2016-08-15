//go:generate reform
package front

//reform:cc_special
type Special struct {
	ID   uint   `reform:"special_id,pk"`
	Name string `reform:"special_name"`
	Pos  int64  `reform:"pos"`
}

//reform:cc_product_special
type ProductSpecial struct {
	ID        uint `reform:"id,pk"`
	SpecialID uint `reform:"special_id"`
	ProductID uint `reform:"product_id"`
}
