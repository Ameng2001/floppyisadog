package oauth

import "github.com/gin-gonic/gin"

func OAuthRegister(router *gin.RouterGroup) {
	router.POST("/tokens", tokensHandler)
	router.POST("/introspect", introspectHandler)
}
