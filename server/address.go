package server

import (
	"net/http"
	"strconv"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostAddr(c *gin.Context) {
	var data front.Address
	if err := c.BindJSON(&data); err != nil {
		front.NewCodeErrv(cerr.InvalidPostBody, err).Abort(c, http.StatusBadRequest)
		return
	}

	err := s.DB.AddressSave(s.TokenUser(c).ID, &data)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, &data)
}

func (s *Server) DeleteAddr(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	err := s.DB.AddressDel(s.TokenUser(c).ID, uint(id))
	if Abort(c, err) {
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
