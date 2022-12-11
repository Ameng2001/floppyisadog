package main

import (
	"context"
	"errors"

	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/utils/crypto"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/jwtverifyserver/tars-protocol/jwtverifyserver"
)

// VerifyImp servant implementation
type VerifyImp struct {
}

// Init servant init
func (imp *VerifyImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destroy
func (imp *VerifyImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *VerifyImp) Verify(ctx context.Context, req *jwtverifyserver.VeifyReq, rsp *jwtverifyserver.VeifyRsp) (int32, error) {
	authz := consts.AuthorizationAnonymousWeb
	rsp.Ret = jwtverifyserver.E_VERIFY_CODE_EVC_SYS_ERR
	if req.Token == "" {
		rsp.Ret = jwtverifyserver.E_VERIFY_CODE_EVC_ERR_TOKEN
		return jwtverifyserver.E_VERIFY_CODE_EVC_ERR_TOKEN, errors.New("bad token")
	}
	uuid, support, err = crypto.RetrieveSessionInformation(req.Token, environment.GetCurrEnv().JWTTokenSecret)
	if err == nil {
		if support {
			authz = consts.AuthorizationSupportUser
		} else {
			authz = consts.AuthorizationAuthenticatedUser
		}

		rsp.Ret = jwtverifyserver.E_VERIFY_CODE_EVC_SUCC
		rsp.Uid = uuid
	} else {
		rsp.Ret = jwtverifyserver.E_VERIFY_CODE_EVC_ERR_TOKEN
	}

	rsp.Context = authz

	return rsp.Ret, nil
}
