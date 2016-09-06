package dbsrv

import (
	"time"

	"gopkg.in/doug-martin/goqu.v3"
	"gopkg.in/reform.v1"

	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
)

func (dbs *DbService) UserVipRebate(tokUsr *models.User, payload *front.VipRebateRequest) error {
	now := time.Now().Unix()
	db := dbs.GetDB()

	// user vip
	ds := dbs.DS.Where(
		goqu.I("$UserID").Eq(tokUsr.ID),
		goqu.I("$ExpiresAt").Gt(now),
		goqu.I("$NotBefore").Lte(now),
	)
	var vip front.VipRebateOrigin
	err := db.DsSelectOneTo(&vip, ds)
	if err != nil && err != reform.ErrNoRows {
		return err
	}

	switch payload.Type {
	case "rebate":
		if vip.ID == 0 {
			return cerr.NotVip
		}
		if vip.Balance == 0 {
			return cerr.VipBalanceEmpty
		}

		// sub vips
		if len(payload.SubIDs) != 2 {
			return cerr.VipRebateSubIDsLen
		}
		var ids []interface{}
		for _, id := range payload.SubIDs {
			if id == 0 {
				return cerr.VipRebateSubIDsHas0
			}
			ids = append(ids, id)
		}
		ds = dbs.DS.Where(
			goqu.I("$User1").Eq(tokUsr.ID),
			goqu.I("$NotBefore").Lte(now),
			goqu.I("$User1Used").Neq(true),
			goqu.I(front.VipRebateOriginTable.PK()).In(ids...),
		)
		ivips, err := db.DsFindAllFrom(front.VipRebateOriginTable, ds)
		if err != nil {
			return err
		}
		if len(ivips) != 2 {
			return cerr.VipRebateSubIDsNoRow
		}

		var orders []uint
		var subsTotal uint
		for _, ivip := range ivips {
			vip := ivip.(*front.VipRebateOrigin)
			subsTotal += vip.Amount
			orders = append(orders, vip.OrderID)
		}

		if subsTotal < vip.Amount*2 {
			return cerr.VipRebateSubTotalSmall
		}

		// use qualifications
		ds := dbs.DS.Where(goqu.I(front.VipRebateOriginTable.PK()).In(ids...))
		_, err = db.DsUpdateColumns(&front.VipRebateOrigin{User1Used: true}, ds, "User1Used")
		if err != nil {
			return err
		}

		// comput and use amount from vip
		var amount uint
		if vip.Amount == vip.Balance {
			amount = vip.Amount / 2
			vip.Balance = vip.Amount - amount
		} else {
			amount = vip.Balance
			vip.Balance = 0
		}
		err = db.UpdateColumns(&vip, "Amount", "Balance")
		if err != nil {
			return err
		}

		// rebate
		err = db.Insert(&front.UserCashRebate{
			UserID:    tokUsr.ID,
			OrderID1:  orders[0],
			OrderID2:  orders[1],
			CreatedAt: now,
			Type:      front.TUserCashRebate,
			Amount:    amount,
			Stages:    dbs.config.Money.UserCashRebateStages,
			DoneAt:    0,
		})
		if err != nil {
			return err
		}

	case "reward":
		// sub vips
		if len(payload.SubIDs) == 0 {
			return cerr.VipRebateSubIDsLen
		}
		var ids []interface{}
		for _, id := range payload.SubIDs {
			if id == 0 {
				return cerr.VipRebateSubIDsHas0
			}
			ids = append(ids, id)
		}
		ds = dbs.DS.Where(
			goqu.I("$User1").Eq(tokUsr.ID),
			goqu.I("$NotBefore").Lte(now),
			goqu.I("$User1Used").Neq(true),
			goqu.I(front.VipRebateOriginTable.PK()).In(ids...),
		)
		ivips, err := db.DsFindAllFrom(front.VipRebateOriginTable, ds)
		if err != nil {
			return err
		}
		if len(ivips) != len(ids) {
			return cerr.VipRebateSubIDsNoRow
		}

		// use qualifications
		ds := dbs.DS.Where(goqu.I(front.VipRebateOriginTable.PK()).In(ids...))
		_, err = db.DsUpdateColumns(&front.VipRebateOrigin{User1Used: true}, ds, "User1Used")
		if err != nil {
			return err
		}

		if vip.ID == 0 {
			// reward in time
			var usr1CashBalance int
			var top front.UserCash
			ds := dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc())
			if err = db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
				return
			}

			for _, ivip := range ivips {
				vip := ivip.(*front.VipRebateOrigin)
				// TODO pause here
				err = db.Insert(&front.UserCash{
					UserID:    order.User1,
					CreatedAt: now,
					Type:      front.TUserCashRebate,
					Amount:    int(toUsr1),
					Balance:   usr1CashBalance,
					OrderID:   order.ID,
				})
				if err != nil {
					return
				}
			}
		} else {
			// freeze reward
		}
	default:
	}
}
