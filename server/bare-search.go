package server

import "github.com/gin-gonic/gin"

func (s *Server) GetNews(c *gin.Context) {
	items, err := s.NewsResource.NewSearcher(c).FindMany()
	ResponseObject(c, items, err)
}
