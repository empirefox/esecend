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

func (s *Server) DeleteCartItems(c *gin.Context) {
	rawids := c.Request.URL.Query()["s"]

	var ids []uint
	if !(len(rawids) == 1 && rawids[0] == "all") {
		for _, rawid := range rawids {
			id, _ := strconv.ParseUint(rawid, 10, 64)
			if id != 0 {
				ids = append(ids, uint(id))
			}
		}
		if len(ids) == 0 {
			front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
			return
		}
	}

	err := s.DB.CartItemDel(s.TokenUser(c).ID, ids)
	if Abort(c, err) {
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
