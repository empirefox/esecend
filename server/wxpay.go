package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) PostWxPayNotify(c *gin.Context) {
	defer c.Request.Body.Close()
	c.XML(http.StatusOK, s.WxClient.OnWxPayNotify(c.Request.Body))
}
