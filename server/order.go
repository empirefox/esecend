package server

import (
	"net/http"
	"strconv"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/lok"
	"github.com/empirefox/esecend/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostCheckout(c *gin.Context) {
	var payload front.CheckoutPayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	data, err := s.DB.CheckoutOrder(s.TokenUser(c), &payload)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) PostOrderPrepay(c *gin.Context) {
	var payload front.OrderPrepayPayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	tokUsr := s.TokenUser(c)
	args, err := s.DB.PrepayOrder(
		tokUsr,
		&payload,
		func(order *front.Order, attach *models.UnifiedOrderAttach) (*front.WxPayArgs, error) {
			return s.WxClient.UnifiedOrder(tokUsr, order, c.ClientIP(), attach)
		},
		s.WxClient.OrderClose,
	)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, args)
}

func (s *Server) PostOrderPay(c *gin.Context) {
	var payload front.OrderPayPayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	order, err := s.DB.PayOrder(s.TokenUser(c), &payload)
	if Abort(c, err) {
		return
	}

	s.DB.GetOrderItems(order)
	c.JSON(http.StatusOK, order)
}

func (s *Server) PostMgrOrderState(c *gin.Context) {
	if !lok.OrderLok.Lock(payload.OrderID) {
		return nil, cerr.OrderTmpLocked
	}
	defer lok.OrderLok.Unlock(payload.OrderID)

	var order front.Order
	err := s.DB.InTx(func(tx *dbsrv.DbService) error {
		return s.DB.MgrOrderState(&order, s.AdminClaims(c), s.WxClient)
	})
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, &order)
}

func (s *Server) PostOrderState(c *gin.Context) {
	var payload front.OrderChangeStatePayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	if !lok.OrderLok.Lock(payload.OrderID) {
		return nil, cerr.OrderTmpLocked
	}
	defer lok.OrderLok.Unlock(payload.OrderID)

	var order front.Order
	err := s.DB.InTx(func(tx *dbsrv.DbService) error {
		return s.DB.OrderChangeState(&order, s.TokenUser(c), &payload, s.WxClient)
	})
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, &order)
}

func (s *Server) GetPaidOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	if !lok.OrderLok.Lock(payload.OrderID) {
		return nil, cerr.OrderTmpLocked
	}
	defer lok.OrderLok.Unlock(payload.OrderID)

	s.DB.InTx(func(dbs *dbsrv.DbService) error {
		order, err := dbs.GetOrder1(s.TokenUser(c), uint(id))
		if Abort(c, err) {
			return err
		}

		if err = s.WxClient.UpdateWxOrderSate(dbs, order); err != nil {
			front.NewCodeErrv(cerr.UpdateWxOrderStateFailed, err).Abort(c, http.StatusInternalServerError)
			return err
		}

		c.JSON(http.StatusOK, order)
		return nil
	})

}

func (s *Server) GetOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	data, err := s.DB.GetOrder1(s.TokenUser(c), uint(id))
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}
