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
)

func (s *Server) Ok(c *gin.Context)       { c.AbortWithStatus(http.StatusOK) }
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
		c.AbortWithError(http.StatusInternalServerError, err)
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return true
	}
	return false
}
