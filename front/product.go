//go:generate reform
package front

import "github.com/empirefox/reform"

//reform:cc_product
type Product struct {
	ID         uint   `reform:"product_id,pk"`
	Name       string `reform:"product_name"`
	Img        string `reform:"img"`
	Intro      string `reform:"intro"`
	Detail     string `reform:"detail"`
	Saled      uint   `reform:"saleCount"`
	CreatedAt  int64  `reform:"create_date"`
	SaledAt    int64  `reform:"time_sale"`
	ShelfOffAt int64  `reform:"time_shelfoff"`
	CategoryID uint   `reform:"cate_id"`
}

//reform:cc_product_sku_att
type ProductAttrId struct {
	ID     uint `reform:"id,pk"`
	SkuID  uint `reform:"sku_id"`
	AttrID uint `reform:"att_id"`
}

type ProductsBundleResponse struct {
	Bundle map[string][]reform.Struct // products
	Skus   []reform.Struct            // Sku
	Attrs  []reform.Struct            // ProductAttrId
}

type ProductsResponse struct {
	Products []reform.Struct // Product
	Skus     []reform.Struct // Sku
	Attrs    []reform.Struct // ProductAttrId
}
