package sms

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/uniuri"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/utils"
	"github.com/gin-gonic/gin"
	"github.com/opensource-conet/alidayu"
	"github.com/patrickmn/go-cache"
)

type LimitedCode struct {
	Code  string
	GenAt int64
}

type Sender struct {
	config    *config.Alidayu
	cache     *cache.Cache
	retryMin  time.Duration
	codeChars []byte
}

func NewSender(config *config.Config, isDebug bool) *Sender {
	dayu := &config.Alidayu
	alidayu.Appkey = dayu.Appkey
	alidayu.AppSecret = dayu.AppSecret
	alidayu.IsDebug = isDebug
	return &Sender{
		config:    dayu,
		cache:     cache.New(dayu.ExpiresMinute*time.Minute, dayu.ClearsMinute*time.Minute),
		retryMin:  dayu.RetryMinSecond * time.Second,
		codeChars: []byte(dayu.CodeChars),
	}
}

// TODO split to server, because we may check if phone is binded already
func (s *Sender) Send(c *gin.Context) {
	phone := c.Param("phone")
	if !utils.RegPhone.MatchString(phone) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if limitedCode, ok := s.cache.Get(phone); ok {
		if time.Now().Add(-s.retryMin).Unix() < limitedCode.(*LimitedCode).GenAt {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
	}

	lcode := LimitedCode{
		Code:  uniuri.NewLenChars(s.config.CodeLen, s.codeChars),
		GenAt: time.Now().Unix(),
	}
	s.cache.Set(phone, &lcode, cache.DefaultExpiration)

	res, err := alidayu.SendOnce(phone, s.config.SignName, s.config.Template, fmt.Sprintf(`{"code":"%s"}`, lcode.Code))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadGateway, res.ResultError)
		c.Abort()
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
