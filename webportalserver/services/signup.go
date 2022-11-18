package services

import (
	"context"
	"net/http"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/gin-gonic/gin"
)

func signUpHandler(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")

	if len(name) <= 0 {
		c.Redirect(http.StatusFound, signUpPath)
	}

	md := make(map[string]string)
	md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService

	acctInfo := new(accountserver.AccountInfo)
	ret, err := outerfactory.Inst().AccountPrx.CreateWithContext(
		context.Background(),
		&accountserver.CreateAccountRequest{Name: name, Email: email},
		acctInfo,
		md,
	)
	if ret != codes.OK {
		logmgr.RERROR("Failed to create account - %v", err)
		c.Redirect(http.StatusFound, signUpPath)
		return
	}

	logmgr.RINFO("New account signup - %v", acctInfo)
	p := configmgr.GetPages().ConfirmPage
	if err = assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
		panic(err)
	}
}
