package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/foauthserver/services/health"
	"github.com/floppyisadog/foauthserver/services/oauth"
	"github.com/floppyisadog/foauthserver/services/user"
	"github.com/floppyisadog/foauthserver/services/web"
	"github.com/floppyisadog/foauthserver/tars-protocol/floppyisadog"
	"github.com/floppyisadog/foauthserver/util"
	"github.com/floppyisadog/foauthserver/util/config"
	"github.com/floppyisadog/foauthserver/util/database"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	sessionName          string = "foauth-session"
	sessionStorePassword string = "floppy-is-a-dog"
)

func main() {
	//01:Get server config
	cfg := tars.GetServerConfig()

	//02: Register Tars protocol servant
	imp := new(foauthImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("foauthImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}

	app := new(floppyisadog.Foauth)
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".FoauthObj")

	//03: Add some customized commands
	//tars.RegisterAdmin("loaddata",cmd.LoadDataHelper)

	//04: Load the server level config file and init database
	tars.AddConfig("foauthserver.conf")
	config.InitConfig(cfg.BasePath + "foauthserver.conf")

	//05: Init database
	db, err := database.InitDB(config.GetConfig())
	if err != nil {
		fmt.Printf("Init database error, err:(%s)\n", err)
		os.Exit(-1)
	}
	defer db.Close()

	//06: Register Http protocol
	mux := &tars.TarsHttpMux{}
	r := gin.Default()
	//r.Use(gin.Recovery())
	//r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//r.LoadHTMLGlob(cfg.BasePath + "template/**/*")
	util.LoadTemplates(cfg.BasePath)
	r.Static("/assets", cfg.BasePath+"assets")

	//07: Registe routes, use cookie as session
	store := cookie.NewStore([]byte(sessionStorePassword))
	r.Use(sessions.Sessions(sessionName, store))

	v1 := r.Group("/v1")
	health.HealthRegister(v1)
	oauth.OAuthRegister(v1.Group("/oauth"))
	user.UserRegister(v1.Group("/user"))
	web.WebRegister(r.Group("/web"))

	mux.Handle("/", r)
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".HttpObj")

	// Run application
	tars.Run()
}
