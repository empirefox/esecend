package server

import (
	"net/http"
	"strconv"

	"github.com/empirefox/reform"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetProfile(c *gin.Context) {
	data, err := s.DB.GetDB().SelectOneFrom(front.ProfileView, "LIMIT 1")
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetTableAll(view reform.View) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := s.DB.GetDB().SelectAllFrom(view, "")
		ResponseArray(c, data, err)
	}
}

func (s *Server) GetEvals(c *gin.Context) {
	productId, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if productId == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	data, err := s.DB.GetDB().FindAllFrom(front.EvalItemView, "$ProductID", productId)
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, data)
}

func (s *Server) GetProductAttrs(c *gin.Context) {
	db := s.DB.GetDB()
	attrs, err := db.SelectAllFrom(front.ProductAttrTable, "")
	if Abort(c, err) {
		return
	}
	grps, err := db.SelectAllFrom(front.ProductAttrGroupTable, "")
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.ProductAttrsResponse{
		Attrs:  attrs,
		Groups: grps,
	})
}

func (s *Server) GetGroupBuy(c *gin.Context) {
	db := s.DB.GetDB()
	items, err := db.SelectAllFrom(front.GroupBuyItemTable, "")
	if Abort(c, err) {
		return
	}

	var args []interface{}
	for _, item := range items {
		args = append(args, item.(*front.GroupBuyItem).SkuID)
	}

	skus, err := db.FindAllFromPK(front.SkuTable, args...)
	if AbortWithoutNoRecord(c, err) {
		return
	}
	c.JSON(http.StatusOK, &front.GroupBuyResponse{
		Items: items,
		Skus:  skus,
	})
}

func (s *Server) GetWishlist(c *gin.Context) {
	db := s.DB.GetDB()
	items, err := db.SelectAllFrom(front.WishItemTable, "")
	if Abort(c, err) {
		return
	}

	var args []interface{}
	for _, item := range items {
		args = append(args, item.(*front.WishItem).ProductID)
	}

	products, err := db.FindAllFromPK(front.ProductTable, args...)
	if AbortWithoutNoRecord(c, err) {
		return
	}
	c.JSON(http.StatusOK, &front.WishListResponse{
		Items:    items,
		Products: products,
	})
}

func (s *Server) GetWallet(c *gin.Context) {
	db := s.DB.GetDB()
	tokUsr := s.TokenUser(c)
	capitalFlows, err := db.FindAllFrom(front.CapitalFlowTable, "$UserID", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}
	pointsList, err := db.FindAllFrom(front.PointsItemTable, "$UserID", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.Wallet{
		CapitalFlows: capitalFlows,
		PointsList:   pointsList,
	})
}

func (s *Server) GetAddrs(c *gin.Context) {
	tokUsr := s.TokenUser(c)
	data, err := s.DB.GetDB().FindAllFrom(front.AddressTable, "$UserID", tokUsr.ID)
	ResponseArray(c, data, err)
}

func (s *Server) GetOrders(c *gin.Context) {
	tokUsr := s.TokenUser(c)
	data, err := s.DB.GetDB().FindAllFrom(front.OrderTable, "$UserID", tokUsr.ID)
	ResponseArray(c, data, err)
}

func (s *Server) GetCart(c *gin.Context) {
	db := s.DB.GetDB()
	tokUsr := s.TokenUser(c)
	items, err := db.FindAllFrom(front.CartItemTable, "$UserID", tokUsr.ID)
	if err == reform.ErrNoRows {
		c.JSON(http.StatusOK, &EmptyArrayJson)
		return
	}
	if Abort(c, err) {
		return
	}

	var args []interface{}
	for _, item := range items {
		args = append(args, item.(*front.CartItem).SkuID)
	}
	skus, err := db.FindAllFromPK(front.SkuTable, args...)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	args = nil
	for _, item := range skus {
		args = append(args, item.(*front.Sku).ProductID)
	}
	products, err := db.FindAllFromPK(front.ProductTable, args...)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.CartResponse{
		Items:    items,
		Skus:     skus,
		Products: products,
	})
}
