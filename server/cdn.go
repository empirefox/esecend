package server

import (
	"net/http"

	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetHeadUptoken(c *gin.Context) {
	c.JSON(http.StatusOK, &front.HeadUptokenResponse{
		HeadToken: s.Cdn.HeadUptoken(s.TokenUser(c).ID),
	})
}
