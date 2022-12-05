package services

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

func loginHandler(c *gin.Context) {
	// if logged in - go away
	// TODO 如何在网关统一设置http鉴权的头（写一个middleware，通过session来记录鉴权）
	if c.GetHeader(consts.AuthorizationHeader) == consts.AuthorizationAuthenticatedUser {
		destination := &url.URL{Host: "myaccount." + environment.GetCurrEnv().ExternalApex, Scheme: "http"}
		c.Redirect(http.StatusFound, destination.String())
		return
	}

	// for GET
	returnTo := c.Query("return_to")

	p := configmgr.GetPages().LoginPage
	p.CsrfField = middleware.TemplateField(c)
	p.ReturnTo = returnTo

	if c.Request.Method == http.MethodPost {
		email := c.PostForm("email")
		password := c.PostForm("password")
		// POST and GET are in the same handler - reset
		returnTo = c.PostForm("return_to")
		// rememberMe=True means that the session is set for a month instead of a day
		rememberMe := len(c.PostForm("remember-me")) > 0

		md := make(map[string]string)
		md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
		acccountInfo := new(accountserver.AccountInfo)
		ret, _ := outerfactory.Inst().AccountPrx.VerifyPasswordWithContext(
			context.Background(),
			&accountserver.VerifyPasswordRequest{Email: email, Password: password},
			acccountInfo,
			md,
		)
		if ret == codes.OK {
			helpers.LoginUser(
				acccountInfo.Uuid,
				acccountInfo.Support,
				rememberMe,
				environment.GetCurrEnv().JWTTokenSecret,
				environment.GetCurrEnv().ExternalApex,
				c.Writer,
			)

			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			outerfactory.Inst().AccountPrx.SyncUserOneWayWithContext(
				context.Background(),
				&accountserver.SyncUserRequest{Uuid: acccountInfo.Uuid},
				md,
			)

			logmgr.RINFO("Logging in user:%s/n", acccountInfo.Email)

			scheme := "https"
			if environment.GetCurrEnv().Name == "development" || environment.GetCurrEnv().Name == "test" {
				scheme = "http"
			}

			if returnTo == "" {
				destination := &url.URL{Host: "app." + environment.GetCurrEnv().ExternalApex, Scheme: scheme}
				returnTo = destination.String()
			} else {
				returnTo = "http://" + returnTo

				// sanitize
				if !isValidSub(returnTo) {
					destination := &url.URL{Host: "myaccount." + environment.GetCurrEnv().ExternalApex, Scheme: scheme}
					returnTo = destination.String()
				}

			}

			c.Redirect(http.StatusFound, returnTo)
			return
		}

		logmgr.RERROR("login attempt denied:(%s)\n", email)
		p.Denied = true
		p.PreviousEmail = email
	}

	if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
		panic(err)
	}
}

// isValidSub returns true if url contains valid subdomain
func isValidSub(sub string) bool {
	u, err := url.Parse(sub)
	if err != nil {
		logmgr.RERROR("can't parse url %v", err)
		return false
	}

	bare := strings.Replace(u.Host, "."+environment.GetCurrEnv().ExternalApex, "", -1)
	for k := range helpers.StaffjoyServices {
		if k == bare {
			return true
		}
	}
	return false
}
