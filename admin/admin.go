package admin

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/front"
	"github.com/mcuadros/go-defaults"
)

type Claims struct {
	jwt.StandardClaims
	OrderID uint `json:"oid"`

	State front.OrderState `json:"state,omitempty"`

	WxRefund   uint `json:"wx_refund,omitempty"`
	CashRefund uint `json:"cash_refund,omitempty"`

	DeliverCom string `json:"deliver_com,omitempty"`
	DeliverNo  string `json:"deliver_no,omitempty"`
}

type Admin struct {
	conf        *config.Security
	GinAdminKey string `default:"ADMIN"`
}

func NewAdmin(conf *config.Config) *Admin {
	a := &Admin{conf: &conf.Security}
	defaults.SetDefaults(a)
	return a
}

func (a *Admin) FindKeyfunc(tok *jwt.Token) (interface{}, error) {
	if tok.Method.Alg() != a.conf.AdminSignType {
		return nil, cerr.InvalidSignAlg
	}

	claims := tok.Claims.(*Claims)
	if claims.ExpiresAt == 0 || claims.ExpiresAt-claims.IssuedAt > 60 {
		return nil, cerr.InvalidTokenExpires
	}
	return []byte(a.conf.AdminKey), nil
}

func (a *Admin) ParseToken(req *http.Request) (*jwt.Token, error) {
	tok, err := request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &Claims{}, a.FindKeyfunc)
	if err != nil {
		fmt.Println("ParseToken", err)
		return nil, cerr.NoAccessToken
	}
	if !tok.Valid {
		return nil, cerr.InvalidAccessToken
	}
	return tok, nil
}
