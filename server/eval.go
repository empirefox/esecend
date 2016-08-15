package server

import (
	"net/http"
	"strconv"

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

	itemId, _ := strconv.ParseUint(c.Query("item"), 10, 64)

	data, err := s.DB.EvalSave(s.TokenUser(c), uint(id), uint(itemId), &payload)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}
