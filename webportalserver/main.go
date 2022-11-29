package main

import (
	"fmt"
	"net/http"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/floppyisadog/appcommon/errorpages"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/services"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
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
	CSRF := csrf.Protect(
		[]byte(configmgr.GetConfig().SigningToken),
		csrf.Domain(configmgr.GetEnvConfig().ExternalApex),
		csrf.Secure(true),
		csrf.Path("/"),
		csrf.CookieName("sjcsrf"),
		csrf.ErrorHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			fmt.Printf("failed CSRF - %s", csrf.FailureReason(req))
			errorpages.Forbidden(res)
		})),
		csrf.FieldName("csrf"),
	)
	r.Use(adapter.Wrap(CSRF))

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
