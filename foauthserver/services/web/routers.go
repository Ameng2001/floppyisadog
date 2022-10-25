package web

import "github.com/gin-gonic/gin"

func WebRegister(router *gin.RouterGroup) {
	router.GET("/register", registerFormHandler)
	router.POST("/register", registerHandler)
	router.GET("/login", loginFormHandler)
	router.POST("/login", loginHandler)
	router.GET("/authorize", authSessionMiddleware(), authorizeFormHandler)
	router.POST("/authorize", authSessionMiddleware(), authorizeHandler)
	router.GET("/logout", authSessionMiddleware(), logoutHandler)
}
