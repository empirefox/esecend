//go:generate reform
package front

//reform:cc_news_item
type NewsItem struct {
	ID        uint   `reform:"id,pk"`
	Icon      string `reform:"icon"`
	Title     string `reform:"title"`
	Detail    string `reform:"detail"`
	CreatedAt int64  `reform:"created_at"`
}
