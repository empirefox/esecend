package dbsrv

import (
	"time"

	"gopkg.in/doug-martin/goqu.v3"

	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/models"
	"github.com/empirefox/esecend/wx"
	"github.com/empirefox/reform"
	"github.com/golang/glog"
)

func (dbs *DbService) UserSetInfo(payload *front.SetUserInfoPayload) error {
	return dbs.GetDB().Update(payload)
}

func (dbs *DbService) FindUserByPhone(phone string) (*models.User, error) {
	usr, err := dbs.GetDB().FindOneFrom(models.UserTable, "$Phone", phone)
	if err != nil {
		return nil, err
	}
	return usr.(*models.User), nil
}

func (dbs *DbService) UserSavePhone(id uint, phone string) (*models.User, error) {
	var usr models.User
	db := dbs.GetDB()
	err := db.FindByPrimaryKeyTo(&usr, id)
	if err != nil {
		return nil, err
	}

	usr.Phone = phone
	if err = db.UpdateColumns(&usr, "Phone"); err != nil {
		return nil, err
	}
	return &usr, nil
}

func (dbs *DbService) UserSetPaykey(id uint, paykey []byte) error {
	data := models.User{
		ID:     id,
		Paykey: &paykey,
	}
	return dbs.GetDB().UpdateColumns(&data, "Paykey")
}

func (dbs *DbService) UserWithdraw(tokUsr *models.User, payload *front.WithdrawPayload) (*front.UserCash, error) {
	if payload.Amount < 100 {
		return nil, cerr.AmountLimit
	}

	db := dbs.GetDB()

	var top front.UserCash
	ds := dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc())
	if err := db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
		return nil, err
	}

	if top.Balance < int(payload.Amount) {
		return nil, cerr.NotEnoughMoney
	}

	now := time.Now().Unix()

	cash := &front.UserCash{
		UserID:    tokUsr.ID,
		CreatedAt: now,
		Type:      front.TUserCashPreWithdraw,
		Amount:    -int(payload.Amount),
		Balance:   top.Balance - int(payload.Amount),
	}
	err := db.Insert(cash)
	if err != nil {
		return nil, err
	}

	data := &wx.TransfersArgs{
		TradeNo: cash.TrackingNumber(),
		OpenID:  tokUsr.OpenId,
		Amount:  payload.Amount,
		Desc:    dbs.config.Money.WithdrawDesc,
		Ip:      payload.Ip,
	}

	result, err := dbs.wc.Transfers(data)
	if err == core.ErrNotFoundSign {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if result["result_code"] != "SUCCESS" {
		glog.Errorln(result["err_code"])
		return nil, cerr.WithdrawFailed
	}

	return cash, nil
}
