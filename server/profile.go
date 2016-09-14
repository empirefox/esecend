package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetMgrReloadProfile(c *gin.Context) {
	err := s.DB.LoadProfile()
	if Abort(c, err) {
		return
	}
	c.AbortWithStatus(http.StatusOK)
}
