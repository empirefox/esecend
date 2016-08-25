package search

import (
	"net/url"

	"github.com/gin-gonic/gin"

	"gopkg.in/doug-martin/goqu.v3"
)

type Pagination struct {
	Total    uint64
	Start    uint64
	Size     uint64
	HasTotal bool
}

type Context struct {
	Resource *Resource

	// user_id must be manual set
	DS *goqu.Dataset

	Pagination Pagination
	Query      url.Values
}

func (res *Resource) NewSearcher(context *gin.Context) *Context {
	c := &Context{
		Resource: res,
		DS:       res.Dbs.DS,
		Query:    context.Request.URL.Query(),
	}
	return c
}

func (res *Resource) NewSearcherFromRaw(rawQuery string) *Context {
	query, _ := url.ParseQuery(rawQuery)
	c := &Context{
		Resource: res,
		DS:       res.Dbs.DS,
		Query:    query,
	}
	return c
}
