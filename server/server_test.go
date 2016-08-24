package server

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/captchar"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/esecend/search"
	"github.com/empirefox/esecend/sec"
	"github.com/empirefox/esecend/sms"
	"github.com/empirefox/esecend/wo2"
	"github.com/empirefox/esecend/wx"
	"github.com/gin-gonic/gin"
)

var (
	isDevMode = true
	log       = logrus.New()

	server *Server
)

func init() {
	server = newServer()
}

func newServer() *Server {
	configFile := os.Getenv("CONFIG")
	if configFile == "" {
		panic("SQL_BASE must be set")
	}

	var err error
	conf, err := config.Load(configFile)
	if err != nil {
		panic(err)
	}

	dbs, err := dbsrv.NewDbService(conf, isDevMode)
	if err != nil {
		panic(err)
	}

	wxClient, err := wx.NewWxClient(conf, dbs)
	if err != nil {
		panic(err)
	}

	captcha, err := captchar.NewCaptchar("../comic.ttf")
	if err != nil {
		panic(err)
	}

	productResource := &search.Resource{
		Conf: conf,
		Dbs:  dbs,
		View: front.ProductTable,
	}
	productResource.SetDefaultFilters()
	productResource.SearchAttrs("Name", "Intro", "Detail")

	secHandler := security.NewHandler(conf, dbs)

	s := &Server{
		IsDevMode:  isDevMode,
		Config:     conf,
		Auther:     wo2.NewAuther(conf, secHandler),
		SecHandler: secHandler,
		Admin:      admin.NewAdmin(conf),
		WxClient:   wxClient,
		SmsSender:  sms.NewSender(conf, isDevMode),
		DB:         dbs,
		Captcha:    captcha,

		ProductResource: productResource,
	}

	auth := s.Auther.Middleware(
		&jwt.Token{
			Claims: &front.TokenClaims{
				OpenId: "open_id",
				UserId: 1,
			},
			Valid: true,
		},
		&models.User{
			ID:     1,
			OpenId: "open_id",
		},
	)
	mustAuthed := s.Auther.MustAuthed

	router := gin.Default()
	a := router.Group("/admin", s.MustAdmin)
	a.POST("/order_state", s.PostMgrOrderState)

	router.POST(s.Config.Security.WxOauthPath, s.Ok)
	router.POST(s.Config.Security.PayNotifyPath, s.PostWxPayNotify)

	router.GET("/profile", s.GetProfile)
	router.GET("/captcha", s.PostCaptcha)
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
	router.POST("/order_state/:id", auth, mustAuthed, s.PostEval)
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

	s.Engine = router
	return s
}
