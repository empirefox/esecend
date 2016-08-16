package server

import (
	"net/http"
	"strconv"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/lok"
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
	payload.Ip = c.ClientIP()

	var args *front.WxPayArgs
	err := s.LockOrderTx(payload.OrderID, func() (cashLocked, pointsLocked bool, err error) {
		args, cashLocked, pointsLocked, err = s.DB.PrepayOrder(s.TokenUser(c), &payload, s.WxClient)
		return
	})
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

	var order front.Order
	err := s.LockOrderTx(claims.OrderID, func() (cashLocked, pointsLocked bool, err error) {
		return s.DB.PayOrder(&order, s.TokenUser(c), &payload)
	})
	if Abort(c, err) {
		return
	}

	s.DB.GetOrderItems(order)
	c.JSON(http.StatusOK, order)
}

func (s *Server) PostMgrOrderState(c *gin.Context) {
	claims := s.AdminClaims(c)

	var order front.Order
	err := s.LockOrderTx(claims.OrderID, func() (cashLocked, pointsLocked bool, err error) {
		return s.DB.MgrOrderState(&order, claims, s.WxClient)
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

	var order front.Order
	err := s.LockOrderTx(payload.ID, func() (cashLocked, pointsLocked bool, err error) {
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

	var order *front.Order
	err := s.LockOrderTx(id, func() (cashLocked, pointsLocked bool, err error) {
		order, err = s.DB.GetBareOrder(s.TokenUser(c), uint(id))
		if err != nil {
			return
		}
		return s.WxClient.UpdateWxOrderSate(dbs, order)
	})
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, order)
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

func (s *Server) LockOrderTx(orderId uint, inTx func() (cashLocked, pointsLocked bool, err error)) (err error) {
	if !lok.OrderLok.Lock(orderId) {
		err = cerr.OrderTmpLocked
		return
	}

	var cashLocked, pointsLocked bool
	defer func() {
		if cashLocked {
			lok.CashLok.Unlock(claims.UserId)
		}
		if pointsLocked {
			lok.PointsLok.Unlock(claims.UserId)
		}
		lok.OrderLok.Unlock(orderId)
	}()

	return s.DB.InTx(func(tx *dbsrv.DbService) (err error) {
		cashLocked, pointsLocked, err = inTx()
		return
	})
}
