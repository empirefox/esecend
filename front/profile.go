//go:generate reform
package front

//reform:cc_profile
type Profile struct {
	Phone            string `reform:"phone"`
	FreeDeliverLine  uint   `reform:"free_delivery_line"`
	DefaultHeadImage string `reform:"default_head_image"`
}

type ProfileResponse struct {
	Profile
	WxAppId     string
	WxScope     string
	WxLoginPath string
}
