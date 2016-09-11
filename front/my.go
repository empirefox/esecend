//go:generate reform
package front

import "github.com/empirefox/reform"

//reform:cc_my_fan
type MyFan struct {
	ID           uint   `reform:"id,pk"`
	CreatedAt    int64  `reform:"create_date"`
	Nickname     string `reform:"name"`
	HeadImageURL string `reform:"avatar"`
	User1        uint   `reform:"parent_id"`
}

type MyFansResponse struct {
	Stores []reform.Struct // Store
	Fans   []reform.Struct // MyFan, alias of models.User
}
