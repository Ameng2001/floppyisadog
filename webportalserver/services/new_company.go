package services

import (
	"context"
	"net/http"
	"net/url"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

const (
	defaultTimezone      = "UTC"
	defaultDayWeekStarts = "monday"
	defaultTeamName      = "Team"
	newCompanyTmpl       = "new_company.tmpl"
	defaultTeamColor     = "744fc6"
)

func newCompanyHandler(c *gin.Context) {
	// if logged in - go away
	authz := helpers.GetAuthFromGinContext(c)
	if authz == consts.AuthorizationAnonymousWeb {
		c.Redirect(http.StatusFound, "/login/")
	}

	if c.Request.Method == http.MethodPost {
		name := c.PostForm("name") // not everything sends this
		timezone := c.PostForm("timezone")
		team := c.PostForm("team")
		if name != "" {
			if timezone == "" {
				timezone = defaultTimezone
			}
			if team == "" {
				team = defaultTeamName
			}

			// fetch current user Infof
			currentUserUUID, err := helpers.GetCurrentUserUUIDFromGinContext(c)
			if err != nil {
				logmgr.RERROR("helpers.GetCurrentUserUUIDFromGinContext return error:(%v)\n", err)
				panic(err)
			}

			md := make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			acccountDTO := new(accountserver.AccountInfo)
			ret, err := outerfactory.Inst().AccountPrx.GetWithContext(
				context.Background(),
				&accountserver.GetAccountRequest{Uuid: currentUserUUID},
				acccountDTO,
				md,
			)
			if ret != codes.OK {
				logmgr.RERROR("AccountPrx.GetWithContext return error:(%d)\n", ret)
				panic(err)
			}

			//create new company
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			companyDTO := new(companyserver.CompanyInfo)
			ret, err = outerfactory.Inst().CompanyPrx.CreateCompanyWithContext(
				context.Background(),
				&companyserver.CreateCompanyRequest{Name: name, Default_timezone: timezone, Default_day_week_starts: defaultDayWeekStarts},
				companyDTO,
				md,
			)
			if ret != codes.OK {
				logmgr.RERROR("CompanyPrx.CreateCompanyWithContext return error:(%v)\n", err)
				panic(err)
			}

			// register current user in directory
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			if ret, err = outerfactory.Inst().CompanyPrx.CreateDirectoryWithContext(
				context.Background(),
				&companyserver.NewDirectoryEntry{Company_uuid: companyDTO.Uuid, Email: acccountDTO.Email},
				new(companyserver.DirectoryEntry),
				md,
			); ret != codes.OK {
				logmgr.RERROR("CompanyPrx.CreateDirectoryWithContext return error:(%v)\n", err)
				panic(err)
			}

			// create admin
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			if ret, err = outerfactory.Inst().CompanyPrx.CreateAdminWithContext(
				context.Background(),
				&companyserver.DirectoryEntryRequest{Company_uuid: companyDTO.Uuid, User_uuid: currentUserUUID},
				new(companyserver.DirectoryEntry),
				md,
			); ret != codes.OK {
				logmgr.RERROR("CompanyPrx.CreateDirectoryWithContext return error:(%v)\n", err)
				panic(err)
			}

			// create team
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			teamDTO := new(companyserver.TeamInfo)
			if ret, err = outerfactory.Inst().CompanyPrx.CreateTeamWithContext(
				context.Background(),
				&companyserver.CreateTeamRequest{Company_uuid: companyDTO.Uuid, Name: team, Color: defaultTeamColor},
				teamDTO,
				md,
			); ret != codes.OK {
				logmgr.RERROR("CompanyPrx.CreateTeamWithContext:(%v)\n", err)
				panic(err)
			}

			// register as worker
			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			if ret, err = outerfactory.Inst().CompanyPrx.CreateWorkerWithContext(
				context.Background(),
				&companyserver.Worker{Company_uuid: companyDTO.Uuid, Team_uuid: teamDTO.Uuid, User_uuid: currentUserUUID},
				new(companyserver.DirectoryEntry),
				md,
			); ret != codes.OK {
				logmgr.RERROR("CompanyPrx.CreateWorkerWithContex:(%v)\n", err)
				panic(err)
			}

			md = make(map[string]string)
			md[consts.AuthorizationMetadata] = consts.AuthorizationWWWService
			outerfactory.Inst().AccountPrx.SyncUserOneWayWithContext(
				context.Background(),
				&accountserver.SyncUserRequest{Uuid: acccountDTO.Uuid},
				md,
			)

			// if environment.GetCurrEnv().Name == "production" && !acccountDTO.Support {
			// 	// TODO:Alert sales of a new account signup
			// 	go func(a *account.Account, c *company.Company) {
			// 		msg := &email.EmailRequest{
			// 			To:       "sales@staffjoy.com",
			// 			Name:     "",
			// 			Subject:  fmt.Sprintf("%s from %s just joined Staffjoy", a.Name, c.Name),
			// 			HtmlBody: fmt.Sprintf("Name: %s<br>Phone: %s<br>Email: %s<br>Company: %s<br>App: https://app.staffjoy.com/#/companies/%s/employees/", a.Name, a.Phonenumber, a.Email, c.Name, c.Uuid),
			// 		}
			// 		mailer, close, err := email.NewClient()
			// 		if err != nil {
			// 			logger.Errorf("unable to initiate email service connection - %s", err)
			// 			return
			// 		}
			// 		defer close()

			// 		ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), md))
			// 		defer cancel()

			// 		if _, err = mailer.Send(ctx, msg); err != nil {
			// 			logger.Errorf("Unable to send email - %s", err)
			// 			return
			// 		}
			// 	}(currentUser, c)
			// }

			// redirect
			logmgr.RINFO("new company signup - %v\n", companyDTO)
			url := url.URL{
				Scheme: "http",
				Host:   "app." + environment.GetCurrEnv().ExternalApex + ":9001",
			}
			c.Redirect(http.StatusFound, url.String())

		}
	}

	p := configmgr.GetPages().NewCompanyPage
	p.CsrfField = middleware.TemplateField(c)
	if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
		panic(err)
	}
}
