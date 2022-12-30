package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/conf"

	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/jwtverifyserver/tars-protocol/jwtverifyserver"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// Init config
	tars.AddConfig("environment.conf")
	tars.AddConfig("jwtverifyserver.conf")
	c, err := conf.NewConf(cfg.BasePath + "accountserver.conf")
	if err != nil {
		//log.RERROR("Parse server config fail", err)
		fmt.Printf("Parse server config fail, err:(%s)\n", err)
	}
	environment.InitEnvironment(c, "/floppyisadog/environment/")

	// New servant imp
	imp := new(VerifyImp)
	err = imp.Init()
	if err != nil {
		fmt.Printf("VerifyImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(jwtverifyserver.Verify)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".VerifyObj")

	// Run application
	tars.Run()
}
