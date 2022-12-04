package repos

import (
	"github.com/floppyisadog/accountserver/managers/logmgr"
	"github.com/floppyisadog/accountserver/models"
	"github.com/floppyisadog/appcommon/utils/database"
)

func FindAccountByEmail(email string) (*models.Account, bool) {
	accountDAO := new(models.Account)
	nofound := database.GetDB().Where("email = LOWER(?)", email).
		First(accountDAO).RecordNotFound()

	if nofound {
		return nil, nofound
	}

	return accountDAO, nofound
}

func FindAccountByPhonenumber(phone string) (*models.Account, bool) {
	accountDAO := new(models.Account)
	nofound := database.GetDB().Where("phone_number = ?", phone).
		First(accountDAO).RecordNotFound()

	if nofound {
		return nil, nofound
	}

	return accountDAO, nofound
}

func FindAccountByUUID(uuid string) (*models.Account, bool) {
	accountDAO := new(models.Account)
	nofound := database.GetDB().Where("id = ?", uuid).
		First(accountDAO).RecordNotFound()

	if nofound {
		return nil, nofound
	}

	return accountDAO, nofound
}

func UpdateAccount(accountDAO *models.Account) bool {
	if err := database.GetDB().Save(accountDAO).Error; err != nil {
		logmgr.RERROR("update account error(%v)\n", err)
		return false
	}

	return true
}

func UpdateAccountFields(uuid string, fields map[string]interface{}) bool {
	query := database.GetDB().Model(new(models.Account)).
		Where("id = ?", uuid)

	if err := query.Update(fields).Error; err != nil {
		logmgr.RERROR("update account error(%s:%v:%v)\n", uuid, fields, err)
		return false
	}

	return true
}
