package server

import (
	"net/http"

	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetNews(c *gin.Context) {
	items, err := s.NewsResource.NewSearcher(c).FindMany()
	ResponseObject(c, items, err)
}

func (s *Server) GetOrders(c *gin.Context) {
	orders, err := s.OrderResource.NewSearcher(c).FindMany()
	if AbortEmptyStructsWithNull(c, orders, err) {
		return
	}

	var ids []interface{}
	for _, order := range orders {
		ids = append(ids, order.(*front.Order).ID)
	}

	items, err := s.DB.GetDB().FindAllFrom(front.OrderItemTable, "$OrderID", ids...)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.OrdersResponse{
		Orders: orders,
		Items:  items,
	})
}
