//go:generate reform
package front

import "time"

//reform:cc_profile
type Profile struct {
	ID                   uint   `reform:"id,pk"` // always 1
	Phone                string `reform:"phone"`
	DefaultHeadImage     string `reform:"default_head_image"`
	UserCashRebateStages uint   `reform:"user_cash_rebate_stages"`
}

type ProfileResponse struct {
	*Profile
	WxAppId     string
	WxScope     string
	WxLoginPath string

	// Config.Order
	EvalTimeoutDay        uint
	CompleteTimeoutDay    uint
	HistoryTimeoutDay     uint
	CheckoutExpiresMinute time.Duration
	WxPayExpiresMinute    time.Duration
	FreeDeliverLine       uint
}
