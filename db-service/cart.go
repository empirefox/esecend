package dbsrv

import (
	"os"
	"time"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/reform"
	"gopkg.in/doug-martin/goqu.v3"
)

// insert or set quantity
func (dbs *DbService) CartItemSave(userId uint, payload *front.SaveToCartPayload) (*front.CartItem, error) {
	if payload.SkuID == 0 {
		return nil, cerr.InvalidSkuId
	}

	data := &front.CartItem{
		UserID:    userId,
		CreatedAt: time.Now().Unix(),
		Name:      payload.Name,
		Img:       payload.Img,
		Type:      payload.Type,
		Price:     payload.Price,
		Quantity:  payload.Quantity,
		SkuID:     payload.SkuID,
	}

	err := dbs.db.InTransaction(func(db *reform.TX) error {
		ds := dbs.DS.Where(goqu.I("$UserID").Eq(userId)).Where(goqu.I("$SkuID").Eq(data.ID))

		// update first
		ra, err := db.DsUpdateStruct(data, ds)
		if err == reform.ErrNoRows {
			return db.Insert(data)
		}
		if err != nil {
			return err
		}
		if ra == 0 {
			return cerr.DbFailed
		}
		table := front.CartItemTable
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

func (dbs *DbService) CartItemDel(userId, id uint) error {
	ds := dbs.DS.Where(goqu.I("$UserID").Eq(userId), goqu.I("$SkuID").Eq(id))
	_, err := dbs.GetDB().DsDelete(front.CartItemTable, ds)
	return err
}
