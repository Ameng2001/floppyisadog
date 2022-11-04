package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"github.com/floppyisadog/botserver/tars-protocol/botserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(BotImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("BotImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(botserver.Bot)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".BotObj")

	// Run application
	tars.Run()
}
