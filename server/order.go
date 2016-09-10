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

	data, err := s.ProductHub.CheckoutOrder(s.TokenUser(c), &payload)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) PostOrderWxPrepay(c *gin.Context) {
	var payload front.OrderPrepayPayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	if payload.OrderID == 0 {
		front.NewCodev(cerr.InvalidPostBody).Abort(c, http.StatusBadRequest)
		return
	}

	tokUsr := s.TokenUser(c)
	_, args, err := s.OrderHub.PrepayOrder(tokUsr.ID, payload.OrderID, c.ClientIP())
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
	if payload.Amount == 0 {
		front.NewCodev(cerr.InvalidPayAmount).Abort(c, http.StatusBadRequest)
		return
	}
	if payload.Key == "" || payload.OrderID == 0 {
		front.NewCodev(cerr.InvalidPostBody).Abort(c, http.StatusBadRequest)
		return
	}

	tokUsr := s.TokenUser(c)
	order, err := s.OrderHub.PayOrder(tokUsr, &payload)
	if Abort(c, err) {
		return
	}

	s.DB.GetOrderItems(&order)
	c.JSON(http.StatusOK, &order)
}

func (s *Server) GetMgrOrderState(c *gin.Context) {
	claims := s.AdminClaims(c)

	var order front.Order
	err := s.OrderHub.MgrOrderState(&order, claims)
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

	tokUsr := s.TokenUser(c)

	var order front.Order
	err := s.OrderHub.OrderChangeState(&order, tokUsr, &payload)
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

	tokUsr := s.TokenUser(c)

	var order *front.Order
	err := s.LockOrderTx(tokUsr.ID, uint(id), func(tx *dbsrv.DbService) (cashLocked, pointsLocked bool, err error) {
		order, err = tx.GetBareOrder(tokUsr, uint(id))
		if err != nil {
			return
		}
		return s.WxClient.UpdateWxOrderSate(tokUsr, tx, order)
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

func (s *Server) LockOrderTx(userId, orderId uint, inTx func(tx *dbsrv.DbService) (cashLocked, pointsLocked bool, err error)) (err error) {
	if !lok.OrderLok.Lock(orderId) {
		err = cerr.OrderTmpLocked
		return
	}

	var cashLocked, pointsLocked bool
	defer func() {
		if cashLocked {
			lok.CashLok.Unlock(userId)
		}
		if pointsLocked {
			lok.PointsLok.Unlock(userId)
		}
		lok.OrderLok.Unlock(orderId)
	}()

	return s.DB.InTx(func(tx *dbsrv.DbService) (err error) {
		cashLocked, pointsLocked, err = inTx(tx)
		return
	})
}
