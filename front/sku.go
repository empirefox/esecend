//go:generate reform
package front

//reform:cc_product_sku
type Sku struct {
	ID          uint   `reform:"sku_id,pk"`
	Stock       uint   `reform:"stock"`
	Img         string `reform:"img"`
	SalePrice   uint   `reform:"sale_price"`
	MarketPrice uint   `reform:"market_price"`
	Freight     uint   `reform:"freight"`
	ProductID   uint   `reform:"product_id"`
	Points      uint   `reform:"points"` // override others
}
