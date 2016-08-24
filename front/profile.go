//go:generate reform
package front

import "time"

//reform:cc_profile
type Profile struct {
	Phone            string `reform:"phone"`
	DefaultHeadImage string `reform:"default_head_image"`
}

type ProfileResponse struct {
	*Profile
	WxAppId     string
	WxScope     string
	WxLoginPath string

	// Config.Order
	EvalTimeoutDay        time.Duration
	CompleteTimeoutDay    time.Duration
	CheckoutExpiresMinute time.Duration
	WxPayExpiresMinute    time.Duration
	Point2Cent            uint
	FreeDeliverLine       uint
}
