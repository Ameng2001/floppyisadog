package repos

import (
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
