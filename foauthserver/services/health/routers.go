package health

import (
	"github.com/gin-gonic/gin"
)

func HealthRegister(router *gin.RouterGroup) {
	router.GET("/health", healthCheck)
}
