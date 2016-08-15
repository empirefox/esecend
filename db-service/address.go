package dbsrv

import (
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
	"gopkg.in/doug-martin/goqu.v3"
)

func (dbs *DbService) AddressSave(userId uint, data *front.Address) error {

	return dbs.db.InTransaction(func(db *reform.TX) error {
		if data.ID == 0 {
			data.UserID = userId
			return db.Insert(data)
		}

		ds := dbs.DS.Where(goqu.I("$UserID").Eq(userId), goqu.I(front.AddressTable.PK()).Eq(data.ID))
		ra, err := db.DsUpdateStruct(data, ds)
		if err != nil {
			return err
		}
		switch ra {
		case 0:
			// TODO validate to insert?
			return reform.ErrNoRows
		case 1:
			return nil
		default:
			return cerr.DbFailed
		}
	})

}

func (dbs *DbService) AddressDel(userId, id uint) error {
	ds := dbs.DS.Where(goqu.I("$UserID").Eq(userId), goqu.I(front.AddressTable.PK()).Eq(id))
	_, err := dbs.GetDB().DsDelete(front.AddressTable, ds)
	return err
}
