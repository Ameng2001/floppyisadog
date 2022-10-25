package user

import "github.com/gin-gonic/gin"

func UserRegister(router *gin.RouterGroup) {
	router.POST("/create", createUserHander)
}
