//go:generate reform
package front

import "github.com/empirefox/reform"

//reform:cc_wishlist
type WishItem struct {
	ID        uint   `reform:"id,pk"`
	UserID    uint   `reform:"user_id"`
	CreatedAt int64  `reform:"created_at"`
	Name      string `reform:"name"`
	Img       string `reform:"img"`
	Price     uint   `reform:"price"`
	ProductID uint   `reform:"product_id"`
}

type WishListResponse struct {
	Items    []reform.Struct // WishItem
	Products []reform.Struct // Product
}

type WishlistSavePayload struct {
	ProductID uint   `binding:"required"`
	Name      string `binding:"required"`
	Img       string `binding:"required"`
	Price     uint   `binding:"required"`
}
