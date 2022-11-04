package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(CompanyImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("CompanyImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(companyserver.Company)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".CompanyObj")

	// Run application
	tars.Run()
}
