package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/floppyisadog/emailserver/tars-protocol/emailserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(EmailImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("EmailImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(emailserver.Email)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".EmailObj")

	// Run application
	tars.Run()
}
