//go:generate reform
package front

import "github.com/empirefox/reform"

//reform:cc_group_buy
type GroupBuyItem struct {
	ID     uint   `reform:"id,pk"`
	Img    string `reform:"img"` // if not present, use sku.Img
	Title  string `reform:"title"`
	Reason string `reform:"reason"`
	Price  uint   `reform:"price"`
	Start  int64  `reform:"start"`
	End    int64  `reform:"end"`
	SkuID  uint   `reform:"sku_id"`
}

type GroupBuyResponse struct {
	Items []reform.Struct // GroupBuyItem
	Skus  []reform.Struct // Sku
}
