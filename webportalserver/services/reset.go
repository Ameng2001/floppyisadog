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
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

func resetHandler(c *gin.Context) {
	p := configmgr.GetPages().ResetPage
	p.CsrfField = middleware.TemplateField(c)
	p.RecaptchaPublic = ""

	if c.Request.Method == http.MethodPost {
		email := c.PostForm("email")
		//TODO:ignore recapcha

		md := make(map[string]string)
		md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
		if ret, err := outerfactory.Inst().AccountPrx.RequestPasswordResetWithContext(
			context.Background(),
			&accountserver.PasswordResetRequest{Email: email},
			md,
		); ret != codes.OK {
			logmgr.RERROR("AccountPrx.RequestPasswordResetWithContext return error:(%v)\n", err)
			panic(err)
		}

		p := configmgr.GetPages().ResetConfirmPage
		if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
			panic(err)
		}
		return
	}

	if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
		panic(err)
	}
}
