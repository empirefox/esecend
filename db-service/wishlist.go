package dbsrv

import (
	"os"
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
)

func (dbs *DbService) WishlistSave(userId uint, payload *front.WishlistSavePayload) (*front.WishItem, error) {
	data := &front.WishItem{
		UserID:    userId,
		CreatedAt: time.Now().Unix(),
		Name:      payload.Name,
		Img:       payload.Img,
		Price:     payload.Price,
		ProductID: payload.ProductID,
	}

	err := dbs.db.InTransaction(func(db *reform.TX) error {
		ds := dbs.DS.Where(goqu.I("$UserID").Eq(userId)).Where(goqu.I("$ProductID").Eq(payload.ProductID))

		// update first
		ra, err := db.DsUpdateStruct(data, ds)
		if err != nil {
			return err
		}
		if ra == 0 {
			return db.Insert(data)
		}
		table := front.WishItemTable
		query, args, err := ds.From(table.Name()).Select(table.PK()).Limit(1).ToSql()
		if err != nil {
			return err
		}
		if err = db.QueryRow(os.Expand(query, table.ToCol), args...).Scan(&data.ID); err != nil {
			return err
		}
		if ra > 1 {
			if _, err = db.DsDelete(table, ds.Where(goqu.I(table.PK()).Neq(data.ID))); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (dbs *DbService) WishlistDel(userId uint, ids []uint) error {
	ds := dbs.DS.Where(goqu.I("$UserID").Eq(userId))
	if len(ids) != 0 {
		ds = ds.Where(goqu.I("$ProductID").In(ids))
	}
	_, err := dbs.GetDB().DsDelete(front.WishItemTable, ds)
	return err
}
