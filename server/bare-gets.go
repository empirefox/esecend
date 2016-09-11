package server

import (
	"net/http"
	"strconv"
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetProfile(c *gin.Context) {
	data, err := s.DB.GetDB().SelectOneFrom(front.ProfileView, "LIMIT 1")
	if Abort(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.ProfileResponse{
		Profile: data.(*front.Profile),

		WxAppId:     s.Config.Weixin.AppId,
		WxLoginPath: s.Config.Security.WxOauthPath,
		//	WxScope     s.Config.Weixin.AppId,

		// Config.Order
		EvalTimeoutDay:        s.Config.Order.EvalTimeoutDay,
		CompleteTimeoutDay:    s.Config.Order.CompleteTimeoutDay,
		HistoryTimeoutDay:     s.Config.Order.HistoryTimeoutDay,
		CheckoutExpiresMinute: s.Config.Order.CheckoutExpiresMinute,
		WxPayExpiresMinute:    s.Config.Order.WxPayExpiresMinute,
		Point2Cent:            s.Config.Order.Point2Cent,
		FreeDeliverLine:       s.Config.Order.FreeDeliverLine,
	})
}

func (s *Server) GetTableAll(view reform.View) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := s.DB.GetDB().SelectAllFrom(view, "")
		ResponseArray(c, data, err)
	}
}

func (s *Server) GetVipIntros(c *gin.Context) {
	now := time.Now().Unix()
	db := s.DB.GetDB()

	ds := dbs.DS.Where(goqu.I("$ExpiresAt").Gt(now), goqu.I("$NotBefore").Lte(now))
	vips, err := db.DsSelectAllFrom(front.VipRebateOriginTable, ds)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	var intros []reform.Struct
	if len(vips) > 0 {
		var ids []interface{}
		for _, vip := range vips {
			ids = append(ids, vip.(*front.VipRebateOrigin).UserID)
		}

		intros, err := db.FindAllFromPK(front.VipIntroTable, ids...)
		if AbortWithoutNoRecord(c, err) {
			return
		}
	}

	c.JSON(http.StatusOK, intros)
}

func (s *Server) GetMyFans(c *gin.Context) {
	db := s.DB.GetDB()
	tokUsr := s.TokenUser(c)

	stores, err := db.FindAllFrom(front.StoreTable, "$User1", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	fans, err := db.FindAllFrom(front.MyFanTable, "$User1", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.MyFansResponse{
		Stores: stores,
		Fans:   fans,
	})
}

// front-end needs to filter not activated
func (s *Server) GetMyQualifications(c *gin.Context) {
	db := s.DB.GetDB()
	tokUsr := s.TokenUser(c)

	vips, err := db.FindAllFrom(front.VipRebateOriginTable, "$User1", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	c.JSON(http.StatusOK, vips)
}

func (s *Server) GetEvals(c *gin.Context) {
	productId, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if productId == 0 {
		front.NewCodev(cerr.InvalidUrlParam).Abort(c, http.StatusBadRequest)
		return
	}

	data, err := s.DB.GetDB().FindAllFrom(front.EvalItemView, front.OrderItemTable.ToCol("ProductID"), productId)
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
	if AbortEmptyStructsWithNull(c, items, err) {
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
	if AbortEmptyStructsWithNull(c, items, err) {
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
	cashes, err := db.FindAllFrom(front.UserCashTable, "$UserID", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	ds := s.DB.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID), goqu.I("$ThawedAt").Eq(0))
	frozen, err := db.DsFindAllFrom(front.UserCashFrozenTable, ds)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	ds = s.DB.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID), goqu.I("$DoneAt").Eq(0))
	rebates, err := db.DsFindAllFrom(front.UserCashRebateTable, ds)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	var rebateItems []reform.Struct
	if len(rebates) != 0 {
		var ids []interface{}
		for _, rebate := range rebates {
			ids = append(ids, rebate.(*front.UserCashRebate).ID)
		}
		rebateItems, err = db.FindAllFromPK(front.UserCashRebateItemTable, ids...)
		if AbortWithoutNoRecord(c, err) {
			return
		}
	}

	points, err := db.FindAllFrom(front.PointsItemTable, "$UserID", tokUsr.ID)
	if AbortWithoutNoRecord(c, err) {
		return
	}

	c.JSON(http.StatusOK, &front.Wallet{
		Cashes:      cashes,
		Frozen:      frozen,
		Rebates:     rebates,
		RebateItems: rebateItems,
		Points:      points,
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
	if AbortEmptyStructsWithNull(c, items, err) {
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

	var products []reform.Struct
	if len(skus) != 0 {
		args = nil
		for _, item := range skus {
			args = append(args, item.(*front.Sku).ProductID)
		}
		products, err = db.FindAllFromPK(front.ProductTable, args...)
		if AbortWithoutNoRecord(c, err) {
			return
		}
	}

	c.JSON(http.StatusOK, &front.CartResponse{
		Items:    items,
		Skus:     skus,
		Products: products,
	})
}
