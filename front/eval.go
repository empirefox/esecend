//go:generate reform
package front

//reform:cc_order_item
type EvalItem struct {
	Eval        string `reform:"comment_content"`
	EvalAt      int64  `reform:"comment_time"` // gen by server
	EvalName    string `reform:"user_name"`    // gen by server
	RateStar    uint   `reform:"starts"`
	RateFit     uint   `reform:"rate_fit"`
	RateServe   uint   `reform:"rate_serve"`
	RateDeliver uint   `reform:"rate_deliver"`
}

type EvalResponse struct {
	Order    *Order // without items
	Evaled   uint
	EvalAt   int64
	EvalName string
}
