package lok

import (
	"sync"
)

var (
	OrderLok  = NewUintLok()
	CashLok   = NewUintLok()
	PointsLok = NewUintLok()
)

type Lok struct {
	l sync.Mutex
	s map[uint]struct{}
}

func NewUintLok() *Lok {
	return &Lok{s: make(map[uint]struct{})}
}

func (ul *Lok) Lock(k uint) bool {
	ul.l.Lock()
	if _, ok := ul.s[k]; ok {
		ul.l.Unlock()
		return false
	}
	ul.s[k] = struct{}{}
	ul.l.Unlock()
	return true
}

func (ul *Lok) Unlock(k uint) {
	ul.l.Lock()
	delete(ul.s, k)
	ul.l.Unlock()
}
