//go:generate reform
package front

import "github.com/dgrijalva/jwt-go"

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
	Privilege string `json:"pvl,omitempty"`
	Phone     string `json:"mob,omitempty"`
	Nonce     string `json:"non,omitempty"`
}

type PreBindPhonePayload struct {
	Phone string
}

type BindPhonePayload struct {
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

type SetPaykeyPayload struct {
	Key  string
	Code string
}
