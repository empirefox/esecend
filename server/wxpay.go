package server

import (
	"net/http"

	"github.com/empirefox/esecend/db-service"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostWxPayNotify(c *gin.Context) {
	defer c.Request.Body.Close()

	res, src := s.WxClient.OnWxPayNotify(c.Request.Body)
	if res != nil {
		c.XML(http.StatusOK, res)
		return
	}

	err := wc.dbs.InTx(func(dbs *dbsrv.DbService) error {
		order, err := dbs.GetBareOrder(nil, id)
		if err != nil {
			return err
		}
		err = dbs.UpdateWxOrderSate(nil, order, src)
		return err
	})

	// must trans to WxReponse
	if err != nil {
		res = NewWxResponse("FAIL", "failed to update trade state")
	} else {
		res = NewWxResponse("SUCCESS", "")
	}

	c.XML(http.StatusOK, res)
}
