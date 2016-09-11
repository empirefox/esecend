//go:generate reform
package models

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/empirefox/esecend/front"
)

var (
	ComparePaykey = bcrypt.CompareHashAndPassword
	EncPaykey     = bcrypt.GenerateFromPassword
)

//reform:cc_member
type User struct {
	// claims below 5 fields
	ID        uint   `reform:"id,pk"`
	OpenId    string `reform:"open_id"`
	Privilege string `reform:"privilege"`
	Phone     string `reform:"phone"`
	User1     uint   `reform:"parent_id"` // TODO add support to login

	CreatedAt int64 `reform:"create_date"`
	UpdatedAt int64 `reform:"update_date"`
	SigninAt  int64 `reform:"last_login"`

	// front.UserInfo
	Nickname     string `reform:"name"`
	Sex          int    `reform:"sex"`
	City         string `reform:"city"`
	Province     string `reform:"province"`
	HeadImageURL string `reform:"avatar"`
	Birthday     int64  `reform:"birthday"`

	Intro string `reform:"intro"`

	UnionId string `reform:"union_id"`

	// for jwt, auto generated when
	// login   sign with new key
	// logout  remove exist keys
	// refresh set old key life with 1min, add the old jti to new head if still valid
	//         sign with new key
	// jwt is saved in mem K-V(jti:key) cache, not in user table
	// Key string

	// RefreshToken is not lookup every time
	// Only query when need refresh
	// Remove when logout
	RefreshToken *[]byte `reform:"refresh_token"` // bcrypt, no expires

	Paykey *[]byte `reform:"paykey"` // for pay, user set, bcrypt
}

func (u *User) Info() *front.UserInfo {
	return &front.UserInfo{
		Nickname:     u.Nickname,
		Sex:          u.Sex,
		City:         u.City,
		Province:     u.Province,
		HeadImageURL: u.HeadImageURL,
		HasPayKey:    u.Paykey != nil && len(*u.Paykey) > 0,
	}
}
