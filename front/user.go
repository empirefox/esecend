//go:generate reform
package front

import "github.com/dgrijalva/jwt-go"

//reform:cc_member
type SetUserInfoPayload struct {
	ID           uint   `reform:"id,pk" json"-"`
	Nickname     string `reform:"name"`
	Sex          int    `reform:"sex"`
	City         string `reform:"city"`
	Province     string `reform:"province"`
	Birthday     int64  `reform:"birthday"`
	CarInsurance string `reform:"car_insurance"`
	InsuranceFee uint   `reform:"insurance_fee"`
	CarIntro     string `reform:"car_intro"`
	Hobby        string `reform:"hobby"`
	Career       string `reform:"career"`
	Demand       string `reform:"demand"`
	Intro        string `reform:"intro"`
	UpdatedAt    int64  `reform:"update_date" json"-"`
}

type UserInfo struct {
	Writable SetUserInfoPayload

	// single modify
	HeadImageURL string

	CreatedAt int64
	UpdatedAt int64

	HasPayKey bool
}

type SetUserInfoResponse struct {
	UpdatedAt int64
}

// out by auth middleware
type UserTokenResponse struct {
	AccessToken  *string
	RefreshToken *string
	User         *UserInfo
}

type TokenClaims struct {
	jwt.StandardClaims
	OpenId string `json:"oid,omitempty"`
	UserId uint   `json:"uid,omitempty"`
	User1  uint   `json:"us1,omitempty"`
	Phone  string `json:"mob,omitempty"`
	Nonce  string `json:"non,omitempty"`
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

type RefreshTokenResponse struct {
	OK          bool
	AccessToken *string
}

type SetPaykeyPayload struct {
	Key  string
	Code string
}
