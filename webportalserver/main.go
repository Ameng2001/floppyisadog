package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/webportalserver/tars-protocol/webportalserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(WebPortalImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("WebPortalImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(webportalserver.WebPortal)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".WebPortalObj")

	// Run application
	tars.Run()
}
