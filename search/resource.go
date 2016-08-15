package search

import (
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/reform"
)

type Resource struct {
	Conf *config.Config
	Dbs  *dbsrv.DbService
	View reform.View

	// used by Context
	scopeMap map[string]*SearchScope

	// output to client
	scopes []*SearchScopeGroup

	// used to overwrite default filters
	filters map[string]*Filter

	SearchHandler func(keyword string, context *Context)
}
