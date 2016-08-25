package server

import (
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"

	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/captchar"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/search"
	"github.com/empirefox/esecend/sec"
	"github.com/empirefox/esecend/sms"
	"github.com/empirefox/esecend/wo2"
	"github.com/empirefox/esecend/wx"
	"github.com/empirefox/gotool/paas"
)

type Server struct {
	*gin.Engine
	IsDevMode  bool
	Config     *config.Config
	Auther     *wo2.Auther
	SecHandler *security.Handler
	Admin      *admin.Admin
	WxClient   *wx.WxClient
	SmsSender  sms.Sender
	DB         *dbsrv.DbService
	Captcha    captchar.Captchar

	ProductResource *search.Resource
}

func (s *Server) BuildEngine() error {

	corsMiddleWare := s.Cors("GET, PUT, POST, DELETE")

	auth := s.Auther.Middleware()
	mustAuthed := s.Auther.MustAuthed

	router := gin.Default()

	router.Use(secure.Secure(secure.Options{
		SSLRedirect: true,
		SSLProxyHeaders: map[string]string{
			"X-Forwarded-Proto": "https",
		},
		IsDevelopment: s.IsDevMode,
	}))
	router.Use(corsMiddleWare)

	router.POST(s.Config.Security.WxOauthPath, s.Ok)
	router.POST(s.Config.Security.PayNotifyPath, s.PostWxPayNotify)

	router.GET("/profile", s.GetProfile)
	router.GET("/carousel", s.GetTableAll(front.CarouselItemTable))
	router.GET("/evals/:product_id", s.GetEvals)
	router.GET("/category", s.GetTableAll(front.CategoryTable))
	router.GET("/product/special/:id", s.GetSpecialProducts)
	router.GET("/product/ls", s.GetProducts)
	router.GET("/product/bundle/:matrix", s.GetProductsBundle)
	router.GET("/product/1/:id", s.GetProduct)
	router.GET("/product/attrs", s.GetProductAttrs)
	router.GET("/groupbuy", s.GetGroupBuy)

	// auth
	router.GET("/refresh_token/:refreshToken", auth, s.HasToken, s.GetRefreshToken)
	router.GET("/captcha", auth, mustAuthed, s.GetCaptcha)
	router.POST("/phone/prebind", auth, mustAuthed, s.PostPreBindPhone)
	router.POST("/phone/bind", auth, mustAuthed, s.PostBindPhone)
	router.GET("/paykey/preset", auth, mustAuthed, s.GetPresetPaykey)
	router.POST("/paykey/set", auth, mustAuthed, s.PostSetPaykey)
	router.GET("/wishlist", auth, mustAuthed, s.GetWishlist)
	router.POST("/wishlist_add", auth, mustAuthed, s.PostWishlistAdd)
	router.DELETE("/wishlist/:id", auth, mustAuthed, s.DeleteWishItem)
	router.GET("/wallet", auth, mustAuthed, s.GetWallet)
	router.GET("/orders", auth, mustAuthed, s.GetOrders)
	router.POST("/checkout", auth, mustAuthed, s.PostCheckout)
	router.POST("/order_pay", auth, mustAuthed, s.PostOrderPay)
	router.POST("/order_wx_pay", auth, mustAuthed, s.PostOrderWxPrepay)
	router.GET("/order/:id", auth, mustAuthed, s.GetOrder)
	router.POST("/order_state", auth, mustAuthed, s.PostOrderState)
	router.GET("/paied_order/:id", auth, mustAuthed, s.GetPaidOrder)
	router.POST("/eval/:id", auth, mustAuthed, s.PostEval)
	router.GET("/cart", auth, mustAuthed, s.GetCart)
	router.POST("/cart", auth, mustAuthed, s.PostCartSave)
	router.DELETE("/cart/:id", auth, mustAuthed, s.DeleteCartItem)
	router.GET("/addrs", auth, mustAuthed, s.GetAddrs)
	router.POST("/addr", auth, mustAuthed, s.PostAddr)
	router.DELETE("/addr/:id", auth, mustAuthed, s.DeleteAddr)
	router.GET("/delivery/:order_id", auth, mustAuthed, s.GetDelivery)
	router.DELETE("/logout", auth, s.DeleteLogout)

	optPathIgnore := make(map[string]bool)
	optPathIgnore[s.Config.Security.PayNotifyPath] = true
	rs := router.Routes()
	for _, r := range rs {
		if r.Method == "OPTIONS" {
			optPathIgnore[r.Path] = true
		}
	}
	for _, r := range rs {
		if !optPathIgnore[r.Path] {
			optPathIgnore[r.Path] = true
			router.OPTIONS(r.Path, s.Ok)
		}
	}

	// for admin
	a := router.Group("/admin", s.MustAdmin)
	a.GET("/order_state", s.GetMgrOrderState)

	s.Engine = router
	return nil
}

func (s *Server) StartRun() error {
	return s.Run(paas.BindAddr)
}
