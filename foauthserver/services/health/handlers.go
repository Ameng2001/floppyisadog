package health

import (
	"net/http"

	"github.com/floppyisadog/appcommon/utils/database"
	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	rows, err := database.GetDB().Raw("SELECT 1=1").Rows()
	healthy := false
	if err == nil {
		healthy = true
	}
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{"healthy": healthy})
}
