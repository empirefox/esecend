package dbsrv

import (
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/doug-martin/goqu.v3"
)

func (s *DbService) SaveProfile(p *front.Profile) error {
	p.ID = 1
	if err := s.GetDB().Update(p); err != nil {
		if err != reform.ErrNoRows {
			return err
		}
		if err = s.insertProfile(p); err != nil {
			return err
		}
	}
	s.SetProfile(p)
	return nil
}

// must be exec after DbService created
func (s *DbService) LoadProfile() error {
	p := new(front.Profile)
	if err := s.GetDB().FindByPrimaryKeyTo(p, 1); err != nil {
		if err != reform.ErrNoRows {
			return err
		}
		defaults.SetDefaults(p)
		if err = s.insertProfile(p); err != nil {
			return err
		}
	}
	s.SetProfile(p)
	return nil
}

func (s *DbService) insertProfile(p *front.Profile) error {
	db := s.GetDB()
	if err := db.Insert(p); err != nil {
		return err
	}
	if p.ID != 1 {
		sql, args, err := s.DS.From(front.ProfileTable.Name()).
			Where(goqu.I(front.ProfileTable.PK()).Eq(p.ID)).
			ToUpdateSql(map[string]interface{}{front.ProfileTable.PK(): 1})
		if err != nil {
			return err
		}
		_, err = db.DsExec(front.ProfileTable, sql, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DbService) SetProfile(p *front.Profile) {
	s.muProfile.Lock()
	s.profile = *p
	s.muProfile.Unlock()
}

func (s *DbService) Profile() front.Profile {
	s.muProfile.RLock()
	defer s.muProfile.RUnlock()
	return s.profile
}
