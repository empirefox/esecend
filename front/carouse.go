//go:generate reform
package front

type BillboardType int

const (
	TBillboardUnknow BillboardType = iota
	TBillboardHome
	TBillboardNews
)

//reform:cc_carousel
type CarouselItem struct {
	ID         uint          `reform:"id,pk"`
	Img        string        `reform:"img"`
	Link       string        `reform:"link"`
	ProductID  uint          `reform:"product_id"`
	CategoryID uint          `reform:"category_id"`
	Pos        int64         `reform:"pos"`
	Billboard  BillboardType `reform:"billboard"`
}
