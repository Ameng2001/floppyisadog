package services

import (
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/gin-gonic/gin"
)

const (
	signUpPath        = "/sign-up/"
	loginPath         = "/login/"
	passwordResetPath = "/password-reset/"
	newCompanyPath    = "/new-company/"
)

func RegisteRoutes(router *gin.RouterGroup) {
	router.GET("/health", healthCheck)
	router.POST("/confirm/", signUpHandler)
	router.GET("/activate/:token", activateGetHandler)
	router.POST("/activate/:token", activatePostHandler)
	router.GET("/reset/:token", confirmResetHandler)
	router.GET(loginPath, loginHandler)
	router.POST(loginPath, loginHandler)
	router.GET("/logout/", logoutHandler)
	router.GET(newCompanyPath, newCompanyHandler)
	router.GET("/breaktime/", breaktimeListHandler)
	router.GET("/breaktime/:slug", breaktimeEpisodeHandler)
	router.GET(passwordResetPath, resetHandler)

	//注册两个api端点
	router.GET("/whoami/", whoamiHandler)
	router.GET("/intercom/", intercomHandler)

	for route, info := range configmgr.GetConfig().Pages.StaticPages {
		router.GET(route, info.Handler)
	}
}
