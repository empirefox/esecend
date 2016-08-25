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
