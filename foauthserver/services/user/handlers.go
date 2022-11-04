package user

import (
	"errors"
	"net/http"

	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/enums"
	"github.com/floppyisadog/foauthserver/services/oauth"
	"github.com/gin-gonic/gin"
)

func createUserHander(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusInternalServerError, codes.NewError("user", err))
		return
	}

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	//check the username and password
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, codes.NewError("user", codes.ErrInvalidUsernameOrPassword))
		return
	}

	//check user existence
	if oauth.UserExists(username) {
		c.JSON(http.StatusBadRequest, codes.NewError("user", errors.New("username taken")))
		return
	}

	//create a new user
	_, err := oauth.CreateUser(
		enums.USER,
		username,
		password,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, codes.NewError("user", err))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
