package server

import (
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"

	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/captchar"
	"github.com/empirefox/esecend/cdn"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/hub"
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
	Cdn        *cdn.Qiniu
	WxClient   *wx.WxClient
	DB         *dbsrv.DbService
	Captcha    captchar.Captchar
	Auther     *wo2.Auther
	SecHandler *security.Handler
	Admin      *admin.Admin
	SmsSender  sms.Sender
	ProductHub *hub.ProductHub
	OrderHub   *hub.OrderHub

	NewsResource    *search.Resource
	ProductResource *search.Resource
	OrderResource   *search.Resource
}

func (s *Server) BuildEngine() {
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

	router.GET("/faketoken", s.GetFakeToken)

	router.POST(s.Config.Security.WxOauthPath, s.Ok)
	router.POST(s.Config.Security.PayNotifyPath, s.PostWxPayNotify)

	router.GET("/profile", s.GetProfile)
	router.GET("/store", s.GetTableAll(front.StoreTable))
	router.GET("/carousel", s.GetTableAll(front.CarouselItemTable))
	router.GET("/evals/:product_id", s.GetEvals)
	router.GET("/category", s.GetTableAll(front.CategoryTable))
	router.GET("/product/ls", s.GetProducts)
	router.GET("/product/bundle/:matrix", s.GetProductsBundle)
	router.GET("/product/1/:id", s.GetProduct)
	router.GET("/product/attrs", s.GetProductAttrs)
	router.GET("/groupbuy", s.GetGroupBuy)
	router.GET("/vips", s.GetVipIntros)
	router.GET("/news", s.GetNews)
	router.GET("/news/1/:id", s.GetNewsItem)

	// auth
	router.GET("/refresh_token/:refreshToken", auth, s.HasToken, s.GetRefreshToken)
	router.GET("/headtoken", auth, mustAuthed, s.GetHeadUptoken)
	router.GET("/captcha", auth, mustAuthed, s.GetCaptcha)
	router.GET("/myfans", auth, mustAuthed, s.GetMyFans)
	router.GET("/myvips", auth, mustAuthed, s.GetMyVips)
	router.GET("/myqualifications", auth, mustAuthed, s.GetMyQualifications)
	router.POST("/rebate", auth, mustAuthed, s.PostUserRebate)
	router.POST("/withdraw", auth, mustAuthed, s.PostUserWithdraw)
	router.POST("/set_user_info", auth, mustAuthed, s.PostSetUserInfo)
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
	router.POST("/checkout_one", auth, mustAuthed, s.PostCheckoutOne)
	router.POST("/order_pay", auth, mustAuthed, s.PostOrderPay)
	router.POST("/order_wx_pay", auth, mustAuthed, s.PostOrderWxPrepay)
	router.GET("/order/:id", auth, mustAuthed, s.GetOrder)
	router.POST("/order_state", auth, mustAuthed, s.PostOrderState)
	router.GET("/paied_order/:id", auth, mustAuthed, s.GetPaidOrder)
	router.POST("/eval/:id", auth, mustAuthed, s.PostEval)
	router.GET("/cart", auth, mustAuthed, s.GetCart)
	router.POST("/cart", auth, mustAuthed, s.PostCartSave)
	router.DELETE("/cart", auth, mustAuthed, s.DeleteCartItems)
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
	a.GET("/reload_profile", s.GetMgrReloadProfile)

	s.Engine = router
}

func (s *Server) StartRun() error {
	go s.ProductHub.Run()
	go s.OrderHub.Run()
	return s.Run(paas.BindAddr)
}
