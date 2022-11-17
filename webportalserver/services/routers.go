package services

import (
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/gin-gonic/gin"
)

func RegisteRoutes(router *gin.RouterGroup) {
	router.GET("/health", healthCheck)
	router.POST("/confirm/", signUpHandler)
	router.GET("/activate/{token}", activateHandler)
	router.GET("/reset/{token}", confirmResetHandler)
	router.GET("/login/", loginHandler)
	router.GET("/logout/", logoutHandler)
	router.GET("/new-company/", newCompanyHandler)
	router.GET("/breaktime/", breaktimeListHandler)
	router.GET("/breaktime/{slug}", breaktimeEpisodeHandler)
	router.GET("/password-reset/", resetHandler)

	for route, info := range configmgr.GetConfig().Pages.StaticPages {
		router.GET(route, info.Handler)
	}
}
