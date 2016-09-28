//go:generate reform
package front

type BillboardType int

const (
	TBillboardUnknow BillboardType = iota
	TBillboardAdSlide
	TBillboardHome
	TBillboardNews
)

//reform:cc_carousel
type CarouselItem struct {
	ID         uint          `reform:"id,pk"`
	Img        string        `reform:"img"`
	Link       string        `reform:"link"`
	ProductID  uint          `reform:"product_id"`
	SpecialID  uint          `reform:"special_id"`
	CategoryID uint          `reform:"category_id"`
	Billboard  BillboardType `reform:"billboard"`
	Pos        int64         `reform:"pos"`
}
