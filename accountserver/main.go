package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(AccountImp)
	err := imp.Init()
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
