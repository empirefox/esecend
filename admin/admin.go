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

type Admin struct {
	conf        *config.Security
	GinAdminKey string `default:"ADMIN"`
}

func (a *Admin) FindAdminKeyfunc(tok *jwt.Token) (interface{}, error) {
	if tok.Method.Alg() != a.conf.AdminSignType {
		return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
	}

	claims := tok.Claims.(*front.TokenClaims)
	if claims.ExpiresAt == 0 || claims.ExpiresAt-claims.IssuedAt > 600 {
		return nil, cerr.InvalidTokenExpires
	}
	return []byte(a.conf.AdminKey), nil
}

func (a *Admin) ParseToken(req *http.Request) (*jwt.Token, error) {
	tok, err := request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &jwt.StandardClaims{}, a.FindAdminKeyfunc)
	if err != nil {
		return nil, cerr.NoAccessToken
	}
	if !tok.Valid {
		return nil, cerr.InvalidAccessToken
	}
	return tok, nil
}
