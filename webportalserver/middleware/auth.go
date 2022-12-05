package middleware

import (
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/helpers"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/gin-gonic/gin"
)

func AuthSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := consts.AuthorizationAnonymousWeb
		uuid, support, err := helpers.GetSession(ctx.Request, environment.GetCurrEnv().JWTTokenSecret)
		if err == nil {
			if support {
				authorization = consts.AuthorizationSupportUser
			} else {
				authorization = consts.AuthorizationAuthenticatedUser
			}
			ctx.Header(consts.CurrentUserHeader, uuid)
		}
		ctx.Header(consts.AuthorizationHeader, authorization)

		ctx.Next()
	}
}
