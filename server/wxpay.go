package server

import (
	"fmt"
	"net/http"

	"github.com/empirefox/esecend/wx"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostWxPayNotify(c *gin.Context) {
	defer c.Request.Body.Close()

	res, src := s.WxClient.OnWxPayNotify(c.Request.Body)
	if res != nil {
		c.XML(http.StatusOK, res)
		return
	}

	var at int64
	var id uint
	_, err := fmt.Sscanf(src["out_trade_no"], "%d-%d", &at, &id)
	if err != nil {
		c.XML(http.StatusOK, wx.NewWxResponse("FAIL", "failed to parse out_trade_no"))
		return
	}

	err = s.OrderHub.OnWxPayNotify(src, id)
	// must trans to WxReponse
	if err != nil {
		res = wx.NewWxResponse("FAIL", "failed to update trade state")
	} else {
		res = wx.NewWxResponse("SUCCESS", "")
	}

	c.XML(http.StatusOK, res)
}
