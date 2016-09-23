package search

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/empirefox/reform"
	"gopkg.in/doug-martin/goqu.v3"
)

//		/:table/ls?typ=n/m
//		&q(query)=bmw
//		&st(start)=100&sz(size)=20&tl(total)=1
//		&sp(scope)=2016style
//		&ft(filter)=Price:gteq:10+Price:lteq:20+Discount:true
//		&ob(order)=Price:desc
func (c *Context) FindMany() ([]reform.Struct, error) {
	hasOrder := c.HandleSearch()
	if c.Pagination.HasTotal && c.Pagination.Total <= c.Pagination.Start {
		return nil, reform.ErrNoRows
	}
	if !hasOrder {
		if table, ok := c.Resource.View.(reform.Table); ok {
			c.DS = c.DS.OrderAppend(goqu.I(table.PK()).Desc())
		}
	}
	return c.Resource.Dbs.GetDB().DsSelectAllFrom(c.Resource.View, c.DS)
}

//		/:table/ls?typ=n/m
//		&q(query)=bmw
//		&st(start)=100&sz(size)=20&tl(total)=1
//		&sp(scope)=2016style
//		&ft(filter)=Price:gteq:10+Price:lteq:20+Discount:true
//		&ob(order)=Price:desc
func (c *Context) HandleSearch() (hasOrder bool) {
	c.HandleQueryScopes()
	c.HandleQueryFilters()
	hasOrder = c.HandleQueryOrder()

	keyword := c.Query.Get("q")
	if keyword != "" && c.Resource.SearchHandler != nil {
		c.Resource.SearchHandler(keyword, c)
	}

	c.HandleQueryPage()

	return
}

// HandleQueryPage handle page from query string:
// &st(start)=100&sz(size)=20&tl(total)=1
func (c *Context) HandleQueryPage() {
	c.Pagination.HasTotal, _ = strconv.ParseBool(c.Query.Get("tl"))
	if c.Pagination.HasTotal {
		var err error
		c.Pagination.Total, err = c.Resource.Dbs.GetDB().DsCount(c.Resource.View, c.DS)
		if err != nil {
			return
		}
	}

	c.Pagination.Start, _ = strconv.ParseUint(c.Query.Get("st"), 10, 64)
	c.Pagination.Size, _ = strconv.ParseUint(c.Query.Get("sz"), 10, 64)
	if c.Pagination.Size == 0 {
		c.Pagination.Size = c.Resource.Conf.Paging.PageSize
	}
	if c.Pagination.Size > c.Resource.Conf.Paging.MaxSize {
		c.Pagination.Size = c.Resource.Conf.Paging.MaxSize
	}

	limit := c.Pagination.Size
	if c.Pagination.HasTotal {
		rest := c.Pagination.Total - c.Pagination.Start
		if rest < limit {
			limit = rest
		}
	}
	c.DS = c.DS.Limit(uint(limit)).Offset(uint(c.Pagination.Start))
}

// HandleQueryOrder handle order from query string:
// &ob(order)=Price:desc
func (c *Context) HandleQueryOrder() (hasOrder bool) {
	orders := strings.Split(c.Query.Get("ob"), ":")
	l := len(orders)
	if l < 1 || l > 2 {
		return
	}

	if col, ok := c.Resource.View.HasCol(orders[0]); ok {
		order := goqu.I(col)
		if l == 2 && orders[1] == "desc" {
			c.DS = c.DS.Order(order.Desc())
		} else {
			c.DS = c.DS.Order(order.Asc())
		}
		return true
	}
	return
}

type searchColKind struct {
	Col  string
	Kind reflect.Kind
}

// SearchAttrs set search attributes, when search resources, will use those columns to search
//     // Search products with its name, code, category's name, brand's name
//	   product.SearchAttrs("Name", "Code", "Category.Name", "Brand.Name")
func (res *Resource) SearchAttrs(fields ...string) {
	if len(fields) == 0 {
		return
	}

	typ := reflect.TypeOf(res.View.NewStruct()).Elem()

	var cks []*searchColKind
	for _, field := range fields {
		if col, ok := res.View.HasCol(field); ok {
			fieldStruct, _ := typ.FieldByName(field)
			indirectType := fieldStruct.Type
			for indirectType.Kind() == reflect.Ptr {
				indirectType = indirectType.Elem()
			}
			cks = append(cks, &searchColKind{Col: col, Kind: indirectType.Kind()})
		}
	}
	tableName := res.View.Name()

	res.SearchHandler = func(keyword string, c *Context) {
		var conditions []string
		var keywords []interface{}

		for _, ck := range cks {

			switch ck.Kind {
			case reflect.String:
				conditions = append(conditions, fmt.Sprintf("upper(%v.%v) like upper(?)", tableName, ck.Col))
				keywords = append(keywords, "%"+keyword+"%")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value, err := strconv.ParseInt(keyword, 10, 64); err == nil {
					conditions = append(conditions, fmt.Sprintf("%v.%v = ?", tableName, ck.Col))
					keywords = append(keywords, value)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if value, err := strconv.ParseUint(keyword, 10, 64); err == nil {
					conditions = append(conditions, fmt.Sprintf("%v.%v = ?", tableName, ck.Col))
					keywords = append(keywords, value)
				}
			case reflect.Float32:
				if value, err := strconv.ParseFloat(keyword, 32); err == nil {
					conditions = append(conditions, fmt.Sprintf("%v.%v = ?", tableName, ck.Col))
					keywords = append(keywords, float32(value))
				}
			case reflect.Float64:
				if value, err := strconv.ParseFloat(keyword, 64); err == nil {
					conditions = append(conditions, fmt.Sprintf("%v.%v = ?", tableName, ck.Col))
					keywords = append(keywords, value)
				}
			case reflect.Bool:
				if value, err := strconv.ParseBool(keyword); err == nil {
					conditions = append(conditions, fmt.Sprintf("%v.%v = ?", tableName, ck.Col))
					keywords = append(keywords, value)
				}
			default:
				conditions = append(conditions, fmt.Sprintf("%v.%v = ?", tableName, ck.Col))
				keywords = append(keywords, keyword)
			}
		}

		// search conditions
		if len(conditions) > 0 {
			c.DS = c.DS.Where(goqu.L(strings.Join(conditions, " OR "), keywords...))
		}
	}
}
