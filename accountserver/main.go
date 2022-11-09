package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/floppyisadog/accountserver/managers/configmgr"
	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/utils/database"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	//Init managers
	tars.AddConfig("accountserver.conf")
	configmgr.InitConfig(cfg.BasePath + "accountserver.conf")

	dbcf := configmgr.GetConfig().Database
	db, err := database.InitDB(dbcf.Type,
		dbcf.User,
		dbcf.Password,
		dbcf.Host,
		dbcf.DatabaseName,
		dbcf.Port,
		dbcf.MaxIdleConns,
		dbcf.MaxOpenConns,
		configmgr.GetConfig().IsDevelopment)
	if err != nil {
		fmt.Printf("Init database error, err:(%s)\n", err)
		os.Exit(-1)
	}
	defer db.Close()

	// New servant imp
	imp := new(AccountImp)
	err = imp.Init()
	if err != nil {
		fmt.Printf("AccountImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(accountserver.Account)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".AccountObj")

	// Run application
	tars.Run()
}
