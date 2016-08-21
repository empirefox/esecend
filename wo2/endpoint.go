package wo2

import (
	"net/http"

	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
)

type Endpoint interface {
	ExchangeToken(code string) (token *oauth2.Token, err error)
	GetUserInfo(accessToken, openId, lang string, httpClient *http.Client) (info *mpoauth2.UserInfo, err error)
}
