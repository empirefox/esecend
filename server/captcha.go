package server

import (
	"net/http"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostCaptcha(c *gin.Context) {
	data, err := s.Captcha.New(s.TokenUser(c).ID)
	if err != nil {
		front.NewCodeErrv(cerr.GenCaptchaFailed, err).Abort(c, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}
