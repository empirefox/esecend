package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostEval(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	var payload front.EvalItem
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	tokUsr := s.TokenUser(c)
	name := []rune(tokUsr.Nickname)
	switch l := len(name); l {
	case 1:
	case 2:
		name[1] = '*'
	case 3:
		name[1] = '*'
		name[2] = '*'
	default:
		for i := 1; i < l-1; i++ {
			name[i] = '*'
		}
	}

	payload.EvalName = string(name)
	payload.EvalAt = time.Now().Unix()

	itemId, _ := strconv.ParseUint(c.Query("item"), 10, 64)

	var order front.Order
	var ra uint
	err := s.OrderHub.EvalSave(&order, &ra, tokUsr, uint(id), uint(itemId), &payload)
	if Abort(c, err) {
		return
	}

	s.DB.GetOrderItems(&order)
	c.JSON(http.StatusOK, &order)
}
