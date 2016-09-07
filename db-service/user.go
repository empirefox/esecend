package dbsrv

import "github.com/empirefox/esecend/models"

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

func (dbs *DbService) UserWithdraw(tokUsr *models.User, payload *front.WithdrawPayload) error {
	if payload.Amount < 100 {
		return cerr.AmountLimit
	}

	var top front.UserCash
	ds = dbs.DS.Where(goqu.I("$UserID").Eq(tokUsr.ID)).Order(goqu.I("$CreatedAt").Desc())
	if err := db.DsSelectOneTo(&top, ds); err != nil && err != reform.ErrNoRows {
		return err
	}

	if top.Balance < payload.Amount {
		return cerr.NotEnoughMoney
	}

	now := time.Now().Unix()
	db := dbs.GetDB()

	cash := &front.UserCash{
		UserID:    tokUsr.ID,
		CreatedAt: now,
		Type:      front.TUserCashPreWithdraw,
		Amount:    -int(payload.Amount),
		Balance:   top.Balance - int(payload.Amount),
	}
	err := db.Insert(cash)
	if err != nil {
		return err
	}

	// TODO split here and back to user goroutine with data

	data := wx.TransfersArgs{
		TradeNo: cash.TrackingNumber(),
		OpenID:  tokUsr.OpenId,
		Amount:  payload.Amount,
		Desc:    dbs.config.Money.WithdrawDesc,
		Ip:      payload.Ip,
	}

	// then next chan
	return nil
}
