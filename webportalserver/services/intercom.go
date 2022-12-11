package services

import (
	"context"
	"net/http"

	"github.com/floppyisadog/accountserver/tars-protocol/accountserver"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/managers/outerfactory"
	"github.com/gin-gonic/gin"
)

type IntercomSettings struct {
	AppID     string `json:"app_id"`
	UserUUID  string `json:"user_id"`
	UserHash  string `json:"user_hash"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

func intercomHandler(c *gin.Context) {
	rsp := IntercomSettings{}
	code := http.StatusOK
	var err error

	rsp.AppID = environment.GetCurrEnv().IntercomAppId
	authz := helpers.GetAuthFromHeader(c)
	switch authz {
	case consts.AuthorizationAnonymousWeb:
		code = http.StatusForbidden
	case consts.AuthorizationSupportUser:
		fallthrough
	case consts.AuthorizationAuthenticatedUser:
		rsp.UserUUID, err = helpers.GetCurrentUserUUIDFromHeader(c.Request.Header)
		if err != nil {
			panic(err)
		}
		rsp.UserHash = crypto.ComputeHmac256(rsp.UserUUID, environment.GetCurrEnv().IntercomSignSecret)
	default:
		logmgr.RERROR("unknown authorization header %s\n", authz)
	}

	if rsp.UserUUID != "" {
		// Get user account info so we can fill in name and email
		md := make(map[string]string)
		md[consts.AuthorizationMetadata] = consts.AuthorizationWhoamiService
		accountDTO := new(accountserver.AccountInfo)
		ret, err := outerfactory.Inst().AccountPrx.GetWithContext(
			context.Background(),
			&accountserver.GetAccountRequest{Uuid: rsp.UserUUID},
			accountDTO,
			md,
		)
		if ret != codes.OK {
			logmgr.RERROR("call AccountPrx.GetWithContext error(%d)(%v)\n", ret, err)
			panic(err)
		}

		rsp.Name = accountDTO.Name
		rsp.Email = accountDTO.Email
		rsp.CreatedAt = accountDTO.Member_since.Seconds

	}

	c.JSON(code, rsp)
}
