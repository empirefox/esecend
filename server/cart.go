package server

import (
	"net/http"
	"strconv"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostCartSave(c *gin.Context) {
	var payload front.SaveToCartPayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	data, err := s.DB.CartItemSave(s.TokenUser(c).ID, &payload)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) DeleteCartItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	err := s.DB.CartItemDel(s.TokenUser(c).ID, uint(id))
	if Abort(c, err) {
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
