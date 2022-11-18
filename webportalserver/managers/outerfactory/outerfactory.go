package outerfactory

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
)

type OuterFactory struct {
	AccountPrx *accountserver.Account
	CompanyPrx *companyserver.Company
}

var instance *OuterFactory

func (c *OuterFactory) initialize() bool {
	comm := tars.NewCommunicator()

	accountObj := configmgr.GetConfig().Outerfactory["AccountObj"]
	c.AccountPrx = new(accountserver.Account)
	comm.StringToProxy(accountObj, c.AccountPrx)

	return true
}

func Inst() *OuterFactory {
	if instance == nil {
		instance = new(OuterFactory)
		instance.initialize()
	}

	return instance
}
