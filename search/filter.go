package search

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"

	"gopkg.in/doug-martin/goqu.v3"
)

type Filter struct {
	Name    string
	Field   string
	Column  string
	Handler FilterHandler
}

type FilterHandler func(params []string, c *Context)

// HandleQueryFilters handle filters from query string:
// &ft(filter)=Price:gteq:10+Price:lteq:20+Discount:true
func (c *Context) HandleQueryFilters() {
	filters := c.Resource.GetFilters()
	for _, ft := range strings.Split(c.Query.Get("ft"), " ") {
		segs := strings.Split(ft, ":")
		if len(segs) > 1 {
			if filter, ok := filters[segs[0]]; ok {
				filter.Handler(segs[1:], c)
			}
		}
	}
}

func (res *Resource) GetFilterNames() (names []string) {
	for name, _ := range res.GetFilters() {
		names = append(names, name)
	}
	return
}

func (res *Resource) GetFilters() map[string]*Filter {
	if res.filters == nil {
		res.SetDefaultFilters()
	}
	return res.filters
}

func (res *Resource) SetDefaultFilters() {
	filters := make(map[string]*Filter)
	for _, field := range res.View.Fields() {
		filter, ok := res.NewFilter(field)
		if ok {
			filters[field] = filter
		} else {
			log.WithFields(logrus.Fields{
				"table": res.View.Name(),
				"field": field,
			}).Errorln("field not found for filter")
		}
	}
	res.filters = filters
}

// NewFilter create Filter from struct field name
func (res *Resource) NewFilter(name string) (*Filter, bool) {
	col, ok := res.View.HasCol(name)
	if !ok {
		return nil, false
	}
	filter := &Filter{
		Name:    name,
		Field:   name,
		Column:  col,
		Handler: res.defaultFilterHandler(col),
	}
	return filter, true
}

// AddFilter add or overwrite the default filter handler
func (res *Resource) AddFilter(name, field string, handler FilterHandler) {
	if res.filters == nil {
		res.filters = make(map[string]*Filter)
	}
	res.filters[name] = &Filter{
		Name:    name,
		Field:   field,
		Column:  res.View.ToCol(field),
		Handler: handler,
	}
}

func (res *Resource) defaultFilterHandler(column string) FilterHandler {
	return func(params []string, c *Context) {
		switch len(params) {
		case 1:
			value, err := strconv.ParseBool(params[0])
			if err == nil {
				c.DS = c.DS.Where(goqu.I(column).Eq(value))
			}

		case 2:
			switch params[0] {
			case "eq":
				c.DS = c.DS.Where(goqu.I(column).Eq(params[1]))
			case "neq":
				c.DS = c.DS.Where(goqu.I(column).Neq(params[1]))
			case "gt":
				c.DS = c.DS.Where(goqu.I(column).Gt(params[1]))
			case "gte":
				c.DS = c.DS.Where(goqu.I(column).Gte(params[1]))
			case "lt":
				c.DS = c.DS.Where(goqu.I(column).Lt(params[1]))
			case "lte":
				c.DS = c.DS.Where(goqu.I(column).Lte(params[1]))
			}
		}
	}
}
