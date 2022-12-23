package services

import (
	"net/http"

	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/gin-gonic/gin"
)

func logoutHandler(c *gin.Context) {
	helpers.Logout(c.Writer, environment.GetCurrEnv().ExternalApex)
	c.Redirect(http.StatusFound, "/")
}
