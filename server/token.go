package server

import (
	"fmt"
	"net/http"

	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetFakeToken(c *gin.Context) {
	ret, err := s.SecHandler.Login(&mpoauth2.UserInfo{
		OpenId:   "OPENID",
		Nickname: "野人",
		UnionId:  "o6_bmasdasdsad6_2sgVt7hMZOPfL",
	}, 0)
	if err != nil {
		fmt.Println(err)
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}
	// front.UserTokenResponse
	c.JSON(200, ret)
}
