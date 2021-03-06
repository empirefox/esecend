package server

import (
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v8"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
	"github.com/gin-gonic/gin"
)

var (
	EmptyArrayJson  = json.RawMessage("[]")
	EmptyObjectJson = json.RawMessage("{}")
	NullJson        = json.RawMessage("null")
)

func (s *Server) Ok(c *gin.Context)       {}
func (s *Server) NotFound(c *gin.Context) { c.AbortWithStatus(http.StatusNotFound) }

func ResponseArray(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, data)
	} else if err == reform.ErrNoRows {
		c.JSON(http.StatusOK, &EmptyArrayJson)
	} else {
		Abort(c, err)
	}
}

func ResponseObject(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, data)
	} else if err == reform.ErrNoRows {
		c.JSON(http.StatusOK, &EmptyObjectJson)
	} else {
		Abort(c, err)
	}
}

func AbortEmptyStructsWithNull(c *gin.Context, data []reform.Struct, err error) bool {
	if Abort(c, err) {
		return true
	}
	if data == nil {
		c.JSON(http.StatusOK, &NullJson)
		return true
	}
	return false
}

func Abort(c *gin.Context, err error) bool {
	if err != nil {
		switch e := err.(type) {
		case cerr.CodedError:
			front.NewCodev(e).Abort(c, http.StatusBadRequest)
			return true
		case validator.ValidationErrors:
			front.NewGinv(err).Abort(c, http.StatusBadRequest)
			return true
		}
		front.NewCodeErrv(cerr.Error, err.Error()).Abort(c, http.StatusInternalServerError)
		return true
	}
	return false
}

func AbortWithoutNoRecord(c *gin.Context, err error) bool {
	if err != nil && err != reform.ErrNoRows {
		switch e := err.(type) {
		case cerr.CodedError:
			front.NewCodev(e).Abort(c, http.StatusBadRequest)
			return true
		case validator.ValidationErrors:
			front.NewGinv(err).Abort(c, http.StatusBadRequest)
			return true
		}
		front.NewCodeErrv(cerr.Error, err.Error()).Abort(c, http.StatusInternalServerError)
		return true
	}
	return false
}
