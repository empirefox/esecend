//go:generate reform
package front

import "github.com/empirefox/reform"

//reform:cc_cart
type CartItem struct {
	ID        uint   `reform:"id,pk"`
	CreatedAt int64  `reform:"created_at"`
	UserID    uint   `reform:"user_id" json:"-"`
	Name      string `reform:"name"`
	Img       string `reform:"img"`
	Type      string `reform:"type"`
	Price     uint   `reform:"price"`
	Quantity  uint   `reform:"quantity"`
	SkuID     uint   `reform:"sku_id"`
}

type SaveToCartPayload struct {
	ID       uint
	Img      string
	Name     string
	Type     string
	Price    uint
	Quantity uint
	SkuID    uint
}

type CartResponse struct {
	Items    []reform.Struct // CartItem
	Skus     []reform.Struct // Sku
	Products []reform.Struct // Product
}
