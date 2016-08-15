//go:generate reform
package front

import "github.com/dgrijalva/jwt-go"

//reform:cc_member_level
type UserLevel struct {
	ID   uint   `reform:"id,pk"`
	Name string `reform:"name"`
}

type UserInfo struct {
	Nickname     string
	Sex          int
	City         string
	Province     string
	HeadImageURL string
	HasPayKey    bool
}

// out by auth middleware
type UserTokenResponse struct {
	AccessToken  *string
	RefreshToken *string
	User         *UserInfo
}

type TokenClaims struct {
	jwt.StandardClaims
	OpenId    string `json:"oid,omitempty"`
	UserId    uint   `json:"uid,omitempty"`
	LevelID   uint   `json:"lid,omitempty"`
	Privilege string `json:"pvl,omitempty"`
	Phone     string `json:"mob,omitempty"`
	Nonce     string `json:"non,omitempty"`
}

type BindPhoneData struct {
	Phone        string `binding:"required"`
	Code         string `binding:"required"`
	CaptchaID    string `binding:"required"`
	Captcha      string `binding:"required"`
	RefreshToken string
}

type BindPhoneResponse struct {
	AccessToken *string
}

type RefreshTokenResponse struct {
	OK          bool
	AccessToken *string
}
