package server

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/gotool/paas"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

func CheckIsSystemMode(c *gin.Context) {
	if paas.IsSystemMode() {
		return
	}
	front.NewCodev(cerr.SystemModeNotAllowed).Abort(c, http.StatusForbidden)
}

func (s *Server) Cors(method string) gin.HandlerFunc {
	return cors.Middleware(cors.Config{
		Origins:         s.Config.Security.Origins,
		Methods:         method,
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          48 * time.Hour,
		Credentials:     false,
		ValidateHeaders: false,
	})
}

func (s *Server) GetRefreshToken(c *gin.Context) {
	token, err := s.SecHandler.RefreshToken(s.Token(c), []byte(c.Param("refreshToken")))

	data := &front.RefreshTokenResponse{OK: true}
	if err == cerr.NoNeedRefreshToken {
		data.OK = false
	} else if Abort(c, err) {
		return
	}

	data.AccessToken = token
	c.JSON(http.StatusOK, data)
}

func (s *Server) DeleteLogout(c *gin.Context) {
	toki, ok := c.Keys[s.Auther.GinJwtKey]
	tok := toki.(*jwt.Token)
	if !ok || !tok.Valid {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	err := s.SecHandler.RevokeToken(tok)
	Abort(c, err)
}

func (s *Server) Claims(c *gin.Context) *front.TokenClaims {
	return s.Token(c).Claims.(*front.TokenClaims)
}

func (s *Server) Token(c *gin.Context) *jwt.Token {
	return s.Auther.GetToken(c)
}

func (s *Server) TokenUser(c *gin.Context) *models.User {
	return c.Keys[s.Auther.GinUserKey].(*models.User)
}

func (s *Server) HasToken(c *gin.Context) {
	_, ok := c.Keys[s.Auther.GinJwtKey]
	if !ok {
		front.NewCodev(cerr.NoAccessToken).Abort(c, http.StatusUnauthorized)
	}
}

func (s *Server) MustAdmin(c *gin.Context) {
	tok, err := s.Admin.ParseToken(c.Request)
	if Abort(c, err) {
		return
	}
	c.Set(s.Admin.GinAdminKey, tok.Claims)
}
