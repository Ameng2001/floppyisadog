package services

import (
	"context"
	"net/http"
	"net/url"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"

	"github.com/floppyisadog/appcommon/errorpages"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

func activateGetHandler(c *gin.Context) {
	page := configmgr.GetPages().ActivatePage
	page.CsrfField = middleware.TemplateField(c)

	token := c.Param("token")
	if len(token) == 0 {
		errorpages.NotFound(c.Writer)
		return
	}

	email, uuid, err := crypto.VerifyEmailConfirmationToken(token, environment.GetCurrEnv().JWTTokenSecret)
	if err != nil {
		c.Redirect(http.StatusFound, passwordResetPath)
	}

	md := make(map[string]string)
	md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService

	acccountInfo := new(accountserver.AccountInfo)
	ret, err := outerfactory.Inst().AccountPrx.GetWithContext(
		context.Background(),
		&accountserver.GetAccountRequest{Uuid: uuid},
		acccountInfo,
		md,
	)
	if ret != codes.OK {
		logmgr.RERROR("call accountserver:get error(%d)(%v)\n", ret, err)
		errorpages.NotFound(c.Writer)
		return
	}

	page.Email = email
	page.Name = acccountInfo.Name
	page.Phonenumber = acccountInfo.Phonenumber

	if err = assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, page.TemplateName, page); err != nil {
		panic(err)
	}
}

func activatePostHandler(c *gin.Context) {
	page := configmgr.GetPages().ActivatePage
	page.CsrfField = middleware.TemplateField(c)

	token := c.Param("token")
	if len(token) == 0 {
		errorpages.NotFound(c.Writer)
		return
	}

	email, uuid, err := crypto.VerifyEmailConfirmationToken(token, environment.GetCurrEnv().JWTTokenSecret)
	if err != nil {
		c.Redirect(http.StatusFound, passwordResetPath)
	}

	md := make(map[string]string)
	md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService

	acccountInfo := new(accountserver.AccountInfo)
	ret, err := outerfactory.Inst().AccountPrx.GetWithContext(
		context.Background(),
		&accountserver.GetAccountRequest{Uuid: uuid},
		acccountInfo,
		md,
	)
	if ret != codes.OK {
		logmgr.RERROR("call accountserver:get error(%d)(%v)\n", ret, err)
		errorpages.NotFound(c.Writer)
		return
	}

	page.Email = email
	page.Name = acccountInfo.Name
	page.Phonenumber = acccountInfo.Phonenumber

	if c.Request.Method == http.MethodPost {
		password := c.PostForm("password")
		name := c.PostForm("name")
		tos := c.PostForm("tos")
		phonenumber := c.PostForm("phonenumber")

		page.Name = name
		page.Phonenumber = phonenumber

		logmgr.RINFO("tos：%v\n", tos)

		if len(password) < 6 {
			page.ErrorMessage = "Your password must be at least 6 characters long"
		}

		if len(tos) <= 0 {
			page.ErrorMessage = "You must agree to the terms and conditions by selecting the checkbox."
		}

		if page.ErrorMessage == "" {
			acccountInfo.Email = email
			acccountInfo.Confirmed_and_active = true
			acccountInfo.Name = name
			acccountInfo.Phonenumber = phonenumber

			md := make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			ret, err := outerfactory.Inst().AccountPrx.UpdateWithContext(
				context.Background(),
				acccountInfo,
				md,
			)
			if ret != codes.OK {
				logmgr.RERROR("call update account error:(%d,%v)\n", ret, err)
				panic(err)
			}

			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			ret, err = outerfactory.Inst().AccountPrx.UpdatePasswordWithContext(
				context.Background(),
				&accountserver.UpdatePasswordRequest{
					Uuid:     acccountInfo.Uuid,
					Password: password,
				},
				md,
			)
			if ret != codes.OK {
				logmgr.RERROR("call updatepassword account error:(%d,%v)\n", ret, err)
				panic(err)
			}

			helpers.LoginUser(
				acccountInfo.Uuid,
				acccountInfo.Support,
				false,
				environment.GetCurrEnv().JWTTokenSecret,
				environment.GetCurrEnv().ExternalApex,
				c.Writer,
			)

			//TODO auditlog
			workerList := new(companyserver.WorkerOfList)
			ret, err = outerfactory.Inst().CompanyPrx.GetWorkerOfWithContext(
				context.Background(),
				&companyserver.WorkerOfRequest{
					User_uuid: acccountInfo.Uuid,
				},
				workerList,
				md,
			)
			if ret != codes.OK {
				logmgr.RERROR("call GetWorkerOfWithContext error:(%d,%v)\n", ret, err)
				panic(err)
			}

			adminList := new(companyserver.AdminOfList)
			ret, err = outerfactory.Inst().CompanyPrx.GetAdminOfWithContext(
				context.Background(),
				&companyserver.AdminOfRequest{
					User_uuid: acccountInfo.Uuid,
				},
				adminList,
				md,
			)
			if ret != codes.OK {
				logmgr.RERROR("call GetAdminOfWithContext error:(%d,%v)\n", ret, err)
				panic(err)
			}

			var destination *url.URL
			if len(adminList.Companies) != 0 || acccountInfo.Support {
				destination = &url.URL{Host: "app." + environment.GetCurrEnv().ExternalApex, Scheme: "http"}
			} else if len(workerList.Teams) != 0 {
				destination = &url.URL{Host: "myaccount." + environment.GetCurrEnv().ExternalApex, Scheme: "http"}
			} else {
				//TODO
				//destination = &url.URL{Host: "www." + environment.GetCurrEnv().ExternalApex, Path: newCompanyPath, Scheme: "http"}
				destination = &url.URL{Host: environment.GetCurrEnv().ExternalApex + ":9001", Path: newCompanyPath, Scheme: "http"}
			}
			c.Redirect(http.StatusFound, destination.String())
		}
	}

	if err = assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, page.TemplateName, page); err != nil {
		panic(err)
	}
}
