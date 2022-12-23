package services

import (
	"context"
	"net/http"
	"net/url"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/errorpages"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

func confirmResetHandler(c *gin.Context) {
	p := configmgr.GetPages().ConfirmResetPage
	p.CsrfField = middleware.TemplateField(c)

	token := c.Param("token")
	if len(token) == 0 {
		errorpages.NotFound(c.Writer)
		return
	}

	email, uuid, err := crypto.VerifyEmailConfirmationToken(token, environment.GetCurrEnv().JWTTokenSecret)
	if err != nil {
		c.Redirect(http.StatusFound, passwordResetPath)
	}

	if c.Request.Method == http.MethodPost {
		// update password
		password := c.PostForm("password")
		if len(password) >= 6 {
			md := make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			accountDTO := new(accountserver.AccountInfo)
			if ret, err := outerfactory.Inst().AccountPrx.GetWithContext(
				context.Background(),
				&accountserver.GetAccountRequest{Uuid: uuid},
				accountDTO,
				md,
			); ret != codes.OK {
				logmgr.RERROR("AccountPrx.GetWithContext return error:(%v)\n", err)
				panic(err)
			}

			accountDTO.Email = email
			accountDTO.Confirmed_and_active = true
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			if ret, err := outerfactory.Inst().AccountPrx.UpdateWithContext(
				context.Background(),
				accountDTO,
				md,
			); ret != codes.OK {
				logmgr.RERROR("AccountPrx.UpdateWithContext return error:(%v)\n", err)
				panic(err)
			}

			// Update password
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			if ret, err := outerfactory.Inst().AccountPrx.UpdatePasswordWithContext(
				context.Background(),
				&accountserver.UpdatePasswordRequest{Uuid: accountDTO.Uuid, Password: password},
				md,
			); ret != codes.OK {
				logmgr.RERROR("AccountPrx.UpdateWithContext return error:(%v)\n", err)
				panic(err)
			}

			// login user
			helpers.LoginUser(
				accountDTO.Uuid,
				accountDTO.Support,
				false,
				environment.GetCurrEnv().JWTTokenSecret,
				environment.GetCurrEnv().ExternalApex,
				c.Writer,
			)

			// Smart redirection - for onboarding purposes
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			workerList := new(companyserver.WorkerOfList)
			if ret, err := outerfactory.Inst().CompanyPrx.GetWorkerOfWithContext(
				context.Background(),
				&companyserver.WorkerOfRequest{
					User_uuid: accountDTO.Uuid,
				},
				workerList,
				md,
			); ret != codes.OK {
				logmgr.RERROR("call GetWorkerOfWithContext error:(%d,%v)\n", ret, err)
				panic(err)
			}

			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			adminList := new(companyserver.AdminOfList)
			if ret, err := outerfactory.Inst().CompanyPrx.GetAdminOfWithContext(
				context.Background(),
				&companyserver.AdminOfRequest{
					User_uuid: accountDTO.Uuid,
				},
				adminList,
				md,
			); ret != codes.OK {
				logmgr.RERROR("call GetAdminOfWithContext error:(%d,%v)\n", ret, err)
				panic(err)
			}

			var destination *url.URL
			if len(adminList.Companies) != 0 || accountDTO.Support {
				destination = &url.URL{Host: "app." + environment.GetCurrEnv().ExternalApex + ":9001", Scheme: "http"}
			} else if len(workerList.Teams) != 0 {
				destination = &url.URL{Host: "myaccount." + environment.GetCurrEnv().ExternalApex + ":9001", Scheme: "http"}
			} else {
				destination = &url.URL{Host: "www." + environment.GetCurrEnv().ExternalApex + ":9001", Path: newCompanyPath, Scheme: "http"}
			}

			c.Redirect(http.StatusFound, destination.String())
		}

		p.ErrorMessage = "Your password must be at least 6 characters long"
	}

	if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
		panic(err)
	}
}
