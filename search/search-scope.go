package search

import "github.com/Sirupsen/logrus"

type SearchScopeGroup struct {
	Name   string
	Scopes []*SearchScope
}

type SearchScope struct {
	Name   string
	Handle func(c *Context)
}

// HandleQueryScopes handle scopes from query string:
// &sp(scope)=2016style
func (c *Context) HandleQueryScopes() {
	if c.Resource.scopeMap != nil {
		if scope, ok := c.Resource.scopeMap[c.Query.Get("sp")]; ok {
			if scope.Handle != nil {
				scope.Handle(c)
			}
		}
	}
}

func (res *Resource) AddSearchScope(name, group string, handle func(*Context)) {
	if res.scopeMap == nil {
		res.scopeMap = make(map[string]*SearchScope)
	}
	if _, ok := res.scopeMap[name]; ok {
		log.WithFields(logrus.Fields{
			"table":       res.View.Name(),
			"serachscope": name,
		}).Fatalln("duplicated scope")
	}
	scope := &SearchScope{
		Name:   name,
		Handle: handle,
	}

	res.scopeMap[name] = scope

	grp := findSearchScope(res.scopes, group)
	if grp == nil {
		grp = &SearchScopeGroup{Name: group}
		res.scopes = append(res.scopes, grp)
	}
	grp.Scopes = append(grp.Scopes, scope)
}

func findSearchScope(groups []*SearchScopeGroup, name string) *SearchScopeGroup {
	for _, group := range groups {
		if group.Name == name {
			return group
		}
	}
	return nil
}
