package server

import (
	"net/http"
	"time"

	"github.com/empirefox/reform"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostBindPhone(c *gin.Context) {
	var data front.BindPhoneData

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

	if !s.Captcha.Validate(data.CaptchaID, data.Captcha) {
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

	usr, err := tx.SaveUserPhone(tokUsr.ID, data.Phone)
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

	var now int64
	refreshToken := []byte(data.RefreshToken)
	if len(refreshToken) > 0 && s.SecHandler.CompareRefreshToken([]byte(usr.RefreshToken), refreshToken) == nil {
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
