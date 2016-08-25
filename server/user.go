package server

import (
	"net/http"
	"time"

	"github.com/empirefox/reform"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/esecend/sms"
	"github.com/empirefox/esecend/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostSetPaykey(c *gin.Context) {
	var payload front.SetPaykeyPayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	tokUsr := s.TokenUser(c)
	if tokUsr.Phone == "" {
		front.NewCodev(cerr.PhoneBindRequired).Abort(c, http.StatusPreconditionFailed)
		return
	}

	if !s.SmsSender.Verify(sms.SetPaykey, tokUsr.ID, tokUsr.Phone, payload.Code) {
		front.NewCodev(cerr.SmsVerifyFailed).Abort(c, http.StatusBadRequest)
		return
	}

	if len(payload.Key) < 6 {
		front.NewCodev(cerr.InvalidPaykey).Abort(c, http.StatusBadRequest)
		return
	}

	paykey, err := models.EncPaykey([]byte(payload.Key), 10)
	if Abort(c, err) {
		return
	}

	err = s.DB.UserSetPaykey(tokUsr.ID, paykey)
	if Abort(c, err) {
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (s *Server) GetPresetPaykey(c *gin.Context) {
	tokUsr := s.TokenUser(c)
	if tokUsr.Phone == "" {
		front.NewCodev(cerr.PhoneBindRequired).Abort(c, http.StatusPreconditionFailed)
		return
	}

	err := s.SmsSender.Send(sms.SetPaykey, tokUsr.ID, tokUsr.Phone)
	if Abort(c, err) {
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (s *Server) PostBindPhone(c *gin.Context) {
	var data front.BindPhonePayload
	if err := c.BindJSON(&data); Abort(c, err) {
		return
	}

	if !utils.RegPhone.MatchString(data.Phone) {
		front.NewCodev(cerr.InvalidPhoneFormat).Abort(c, http.StatusBadRequest)
		return
	}

	tokUsr := s.TokenUser(c)
	if tokUsr.Phone == data.Phone {
		front.NewCodev(cerr.RebindSamePhone).Abort(c, http.StatusBadRequest)
		return
	}

	if !s.SmsSender.Verify(sms.BindPhone, tokUsr.ID, data.Phone, data.Code) {
		front.NewCodev(cerr.SmsVerifyFailed).Abort(c, http.StatusBadRequest)
		return
	}

	if !s.Captcha.Verify(tokUsr.ID, data.CaptchaID, data.Captcha) {
		front.NewCodev(cerr.CaptchaRejected).Abort(c, http.StatusBadRequest)
		return
	}

	tx, err := s.DB.Tx()
	if err != nil {
		front.NewCodeErrv(cerr.DbFailed, err).Abort(c, http.StatusInternalServerError)
		return
	}
	defer tx.RollbackIfNeeded()

	_, err = tx.FindUserByPhone(data.Phone)
	if err == nil {
		front.NewCodev(cerr.PhoneOccupied).Abort(c, http.StatusBadRequest)
		return
	}
	if err != reform.ErrNoRows {
		front.NewCodeErrv(cerr.DbFailed, err).Abort(c, http.StatusInternalServerError)
		return
	}

	usr, err := tx.UserSavePhone(tokUsr.ID, data.Phone)
	if err == reform.ErrNoRows {
		front.NewCodeErrv(cerr.UserNotFound, err).Abort(c, http.StatusInternalServerError)
		return
	}

	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		front.NewCodeErrv(cerr.DbFailed, err).Abort(c, http.StatusInternalServerError)
		return
	}

	if usr.RefreshToken == nil || len(*usr.RefreshToken) == 0 {
		c.JSON(http.StatusOK, &EmptyObjectJson)
		return
	}

	var now int64
	refreshToken := []byte(data.RefreshToken)
	if len(refreshToken) > 0 && s.SecHandler.CompareRefreshToken(*usr.RefreshToken, refreshToken) == nil {
		now = time.Now().Unix()
	} else {
		now = s.Claims(c).IssuedAt
	}
	tok, err := s.SecHandler.NewTokenWithIat(usr, now)
	if err != nil {
		c.JSON(http.StatusOK, &EmptyObjectJson)
		return
	}

	c.JSON(http.StatusOK, &front.RefreshTokenResponse{OK: true, AccessToken: tok})
}

func (s *Server) PostPreBindPhone(c *gin.Context) {
	var payload front.PreBindPhonePayload
	if err := c.BindJSON(&payload); Abort(c, err) {
		return
	}

	if !utils.RegPhone.MatchString(payload.Phone) {
		front.NewCodev(cerr.InvalidPostBody).Abort(c, http.StatusBadRequest)
		return
	}

	tokUsr := s.TokenUser(c)
	if tokUsr.Phone == payload.Phone {
		front.NewCodev(cerr.RebindSamePhone).Abort(c, http.StatusBadRequest)
		return
	}

	_, err := s.DB.FindUserByPhone(payload.Phone)
	if err == nil {
		front.NewCodev(cerr.PhoneOccupied).Abort(c, http.StatusBadRequest)
		return
	}
	if err != reform.ErrNoRows {
		front.NewCodeErrv(cerr.DbFailed, err).Abort(c, http.StatusInternalServerError)
		return
	}

	err = s.SmsSender.Send(sms.BindPhone, tokUsr.ID, payload.Phone)
	if Abort(c, err) {
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
