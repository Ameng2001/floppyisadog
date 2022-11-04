package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/floppyisadog/smsserver/tars-protocol/smsserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(SmsImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("SmsImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(smsserver.Sms)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".SmsObj")

	// Run application
	tars.Run()
}
