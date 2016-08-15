package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/delivery"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetDelivery(c *gin.Context) {
	orderId, _ := strconv.ParseUint(c.Param("order_id"), 10, 64)
	if orderId == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	order, err := s.DB.GetOrder1(s.TokenUser(c), uint(orderId))
	if Abort(c, err) {
		return
	}

	res, err := delivery.QueryRemote(order.DeliverCom, order.DeliverNo)
	if err != nil {
		front.NewCodeErrv(cerr.RemoteHTTPFailed, err).Abort(c, http.StatusBadGateway)
		return
	}

	c.Status(http.StatusOK)
	defer res.Body.Close()
	io.Copy(c.Writer, res.Body)
}
