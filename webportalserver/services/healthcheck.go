package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	healthy := true
	c.JSON(http.StatusOK, gin.H{"healthy": healthy})
}
