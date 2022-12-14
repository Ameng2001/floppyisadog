package services

import (
	"context"
	"net/http"

	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/companyserver/tars-protocol/companyserver"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/gin-gonic/gin"
)

type IAm struct {
	Support  bool   `json:"support"`
	UserUUID string `json:"user_uuid"`
	//Token    string                      `json:"token"`
	Authz  string                      `json:"authz"`
	Worker *companyserver.WorkerOfList `json:"worker"`
	Admin  *companyserver.AdminOfList  `json:"admin"`
}

func whoamiHandler(c *gin.Context) {
	rsp := IAm{}
	code := http.StatusOK
	var err error

	authz := helpers.GetAuthFromGinContext(c)
	switch authz {
	case consts.AuthorizationAnonymousWeb:
		code = http.StatusForbidden
	case consts.AuthorizationAuthenticatedUser:
		rsp.UserUUID, err = helpers.GetCurrentUserUUIDFromGinContext(c)
		if err != nil {
			panic(err)
		}
	case consts.AuthorizationSupportUser:
		rsp.Support = true
		rsp.UserUUID, err = helpers.GetCurrentUserUUIDFromGinContext(c)
		if err != nil {
			panic(err)
		}
	default:
		logmgr.RERROR("unknown authorization header %s\n", authz)
	}
	rsp.Authz = authz

	if rsp.UserUUID != "" {
		//rsp.Token = helpers.GetStoredJWT(c.Request)
		// Get worker stuff
		md := make(map[string]string)
		md[consts.AuthorizationMetadata] = consts.AuthorizationWhoamiService

		rsp.Worker = new(companyserver.WorkerOfList)
		ret, err := outerfactory.Inst().CompanyPrx.GetWorkerOfWithContext(
			context.Background(),
			&companyserver.WorkerOfRequest{User_uuid: rsp.UserUUID},
			rsp.Worker,
			md,
		)
		if ret != codes.OK {
			logmgr.RERROR("call CompanyPrx.GetWorkerOfWithContext error(%d)(%v)\n", ret, err)
			panic(err)
		}

		md = make(map[string]string)
		md[consts.AuthorizationMetadata] = consts.AuthorizationWhoamiService
		rsp.Admin = new(companyserver.AdminOfList)
		ret, err = outerfactory.Inst().CompanyPrx.GetAdminOfWithContext(
			context.Background(),
			&companyserver.AdminOfRequest{User_uuid: rsp.UserUUID},
			rsp.Admin,
			md,
		)
		if ret != codes.OK {
			logmgr.RERROR("call CompanyPrx.GetAdminOfWithContex error(%d)(%v)\n", ret, err)
			panic(err)
		}
	}

	c.JSON(code, rsp)
}
