package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/appcommon/errorpages"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/floppyisadog/webportalserver/services"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := tars.GetServerConfig()

	tars.AddConfig("webportalserver.conf")
	configmgr.InitConfig(cfg.BasePath + "webportalserver.conf")

	//Register Http protocol
	mux := &tars.TarsHttpMux{}
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//CSRF middleware
	middleware.InitCSRF(r, configmgr.GetConfig().SigningToken)

	//load assets and templates
	errorpages.InitAssets()
	assetsmgr.LoadAssets()
	assetsmgr.RegisteStatic(r)

	//registe routers
	root := r.Group("/")
	services.RegisteRoutes(root)
	r.NoRoute(func(ctx *gin.Context) {
		errorpages.NotFound(ctx.Writer)
	})

	//Add http servant
	mux.Handle("/", r)
	tars.AddHttpServant(mux, cfg.App+"."+cfg.Server+".HttpObj")

	// Run application
	tars.Run()
}
