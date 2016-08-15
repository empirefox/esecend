package admin

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/front"
)

type Claims struct {
	jwt.StandardClaims
	UserId  uint `json:"uid"`
	OrderID uint `json:"oid"`

	State front.TradeState `json:"state"`

	WxRefund     uint `json:"wx_refund,omitempty"`
	CashRefund   uint `json:"cash_refund,omitempty"`
	PointsRefund uint `json:"points_refund,omitempty"`

	DeliverCom string `json:"deliver_com,omitempty"`
	DeliverNo  string `json:"deliver_no,omitempty"`
}

type Admin struct {
	conf        *config.Security
	GinAdminKey string `default:"ADMIN"`
}

func (a *Admin) FindAdminKeyfunc(tok *jwt.Token) (interface{}, error) {
	if tok.Method.Alg() != a.conf.AdminSignType {
		return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
	}

	claims := tok.Claims.(*Claims)
	if claims.ExpiresAt == 0 || claims.ExpiresAt-claims.IssuedAt > 30 {
		return nil, cerr.InvalidTokenExpires
	}
	return []byte(a.conf.AdminKey), nil
}

func (a *Admin) ParseToken(req *http.Request) (*jwt.Token, error) {
	tok, err := request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &Claims{}, a.FindAdminKeyfunc)
	if err != nil {
		return nil, cerr.NoAccessToken
	}
	if !tok.Valid {
		return nil, cerr.InvalidAccessToken
	}
	return tok, nil
}
