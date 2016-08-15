package wo2

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/buger/jsonparser"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/empirefox/esecend/config"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
)

var log = logrus.New()

type SecurityHandler interface {
	Login(userinfo *mpoauth2.UserInfo) (ret interface{}, err error)
	ParseToken(c *gin.Context) (tok *jwt.Token, tokUsr interface{}, err error)
}

type Auther struct {
	GinJwtKey  string `default:"claims"`
	GinUserKey string `default:"tokUser"`

	wx          *config.Weixin
	wxOauthPath string
	secHandler  SecurityHandler
	endpoint    *mpoauth2.Endpoint
}

func NewAuther(conf *config.Config, secHandler SecurityHandler) *Auther {
	return &Auther{
		wxOauthPath: conf.Security.WxOauthPath,
		wx:          &conf.Weixin,
		secHandler:  secHandler,
	}
}

// Middleware proccess Login related logic.
// It does not block the user handler, just try to retrieve Token.Claims.
func (a *Auther) Middleware(iuser ...interface{}) gin.HandlerFunc {
	if a == nil {
		panic("goauth Auther is nil")
	}
	a.loadDefault()

	if len(iuser) > 0 {
		return func(c *gin.Context) {
			c.Set(a.GinUserKey, iuser[0])
		}
	}

	return func(c *gin.Context) {
		if c.Request.URL.Path == a.wxOauthPath && c.Request.Method == "POST" {
			if err := a.authHandle(c); err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
		} else {
			tok, user, err := a.secHandler.ParseToken(c)
			if err == nil {
				c.Set(a.GinJwtKey, tok)
				c.Set(a.GinUserKey, user)
			}
		}
	}
}

func (a *Auther) MustAuthed(c *gin.Context) {
	tok, ok := c.Keys[a.GinJwtKey]
	if !ok || !tok.(*jwt.Token).Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (a *Auther) loadDefault() {
	if result, err := govalidator.ValidateStruct(a); !result {
		panic(err)
	}

	defaults.SetDefaults(a)
	a.endpoint = mpoauth2.NewEndpoint(a.wx.AppId, a.wx.ApiKey)
}

func (a *Auther) authHandle(c *gin.Context) error {
	raw, _ := ioutil.ReadAll(c.Request.Body)
	log.Debugf("Code Body:%s\n", raw)
	code, err := jsonparser.GetUnsafeString(raw, "code")
	if err != nil {
		return err
	}

	client := &oauth2.Client{Endpoint: a.endpoint}
	tok, err := client.ExchangeToken(code)
	if err != nil {
		return err
	}

	userinfo, err := mpoauth2.GetUserInfo(tok.AccessToken, tok.OpenId, "", nil)
	if err != nil {
		return err
	}

	ret, err := a.secHandler.Login(userinfo)
	if err != nil {
		return err
	}
	// front.UserTokenResponse
	c.JSON(200, ret)
	c.Abort()
	return nil
}

func (a *Auther) GetToken(c *gin.Context) *jwt.Token {
	return c.Keys[a.GinJwtKey].(*jwt.Token)
}
