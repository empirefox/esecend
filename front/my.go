//go:generate reform
package front

//reform:cc_member
type MyFan struct {
	ID           uint   `reform:"id,pk"`
	CreatedAt    int64  `reform:"create_date"`
	Nickname     string `reform:"name"`
	HeadImageURL string `reform:"avatar"`
	User1        uint   `reform:"parent_id"`
}
