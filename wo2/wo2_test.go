package wo2

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/l"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Level = logrus.DebugLevel
}

func TestAuther_Oauth(t *testing.T) {
	auther := newAuther(newSecHandler())
	res := request(auther, "POST", "/auth", strings.NewReader(`{"code":"CODE"}`))
	require.Equal(t, 200, res.Code)
	require.Equal(t, `{"token":"TOKEN"}`, strings.Trim(res.Body.String()))
}

func request(auther *Auther, method, path string, payload io.Reader) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Use(auther.Middleware())
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, payload)
	r.ServeHTTP(res, req)
	return res
}

func newSecHandler() *secHandler {
	return &secHandler{
		loginReturn: gin.H{"token": "TOKEN"},
		loginErr:    nil,
		parsedToken: new(jwt.Token),
		parsedUser:  gin.H{"ID": 111},
		parsedErr:   nil,
	}
}

type secHandler struct {
	loginReturn interface{}
	loginErr    error
	parsedToken *jwt.Token
	parsedUser  interface{}
	parsedErr   error
}

func (h *secHandler) Login(userinfo *mpoauth2.UserInfo) (ret interface{}, err error) {
	return h.loginReturn, h.loginErr
}

func (h *secHandler) ParseToken(c *gin.Context) (tok *jwt.Token, tokUsr interface{}, err error) {
	return h.parsedToken, h.parsedUser, h.parsedErr
}

type oauth2HttpClientTransport struct{}

func (t oauth2HttpClientTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")
	responseBody := ""
	if req.URL.Path == "/sns/oauth2/access_token" {
		responseBody = `
{
  "access_token":"ACCESS_TOKEN",
  "expires_in":7200,
  "refresh_token":"REFRESH_TOKEN",
  "openid":"OPENID",
  "scope":"SCOPE",
  "unionid":"o6_bmasdasdsad6_2sgVt7hMZOPfL"
}`
	} else {
		responseBody = `
{
  "access_token":"ACCESS_TOKEN",
  "expires_in":7200,
  "refresh_token":"REFRESH_TOKEN",
  "openid":"OPENID",
  "scope":"SCOPE"
}`
	}
	log.WithFields(l.Locate(logrus.Fields{
		"path": req.URL.Path,
	})).Debugf("requested")
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

type getUserInfoHttpClientTransport struct{}

func (t getUserInfoHttpClientTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")
	responseBody := `
{
  "openid":"OPENID",
  "nickname":"NICKNAME",
  "sex":1,
  "province":"PROVINCE",
  "city":"CITY",
  "country":"COUNTRY",
  "headimgurl":"http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/46", 
	"privilege":[
	"PRIVILEGE1",
	"PRIVILEGE2"
   ],
   "unionid":"o6_bmasdasdsad6_2sgVt7hMZOPfL"
}`
	log.WithFields(l.Locate(logrus.Fields{
		"path": req.URL.Path,
	})).Debugf("requested")
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

func newAuther(h *secHandler) *Auther {
	oauth2HttpClient := &http.Client{Transport: new(oauth2HttpClientTransport)}
	getUserInfoHttpClient := &http.Client{Transport: new(getUserInfoHttpClientTransport)}
	return &Auther{

		Oauth2HttpClient:      oauth2HttpClient,
		GetUserInfoHttpClient: getUserInfoHttpClient,

		wx:          new(config.Weixin),
		wxOauthPath: "/auth",
		secHandler:  h,
	}
}
