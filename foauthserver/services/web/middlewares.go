package web

import (
	"encoding/gob"

	"github.com/floppyisadog/foauthserver/services/oauth"
	"github.com/floppyisadog/foauthserver/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserSession has user data stored in a session after logging in
type UserSession struct {
	ClientID     string
	Username     string
	AccessToken  string
	RefreshToken string
}

var (
	USER_SESSION_KEY string = "user-session"
)

func init() {
	// Register a new datatype for storage in sessions
	gob.Register(new(UserSession))
}

func authSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionUser := session.Get(USER_SESSION_KEY).(*UserSession)
		if sessionUser == nil {
			query := ctx.Request.URL.Query()
			query.Set("login_redirect_uri", ctx.Request.URL.Path)
			util.RedirectWithQueryString("/login", query, ctx)
			return
		}

		if err := authenticateWithSession(sessionUser); err != nil {
			query := ctx.Request.URL.Query()
			query.Set("login_redirect_uri", ctx.Request.URL.Path)
			util.RedirectWithQueryString("/login", query, ctx)
			return
		}

		// update user session
		session.Set(USER_SESSION_KEY, sessionUser)
		session.Save()
	}
}

func authenticateWithSession(user *UserSession) error {
	// Try to authenticate with the stored access token
	_, err := oauth.ValidateAccessToken(user.AccessToken)
	if err == nil {
		//AccessToken valid
		return nil
	}

	// Access token might be expired, let's try refreshing...
	client, err := oauth.FindClientByClientKey(user.ClientID)
	if err != nil {
		return nil
	}

	// Validate the refreshToken
	sessionRefreshToken, err := oauth.ValidateRefreshToken(user.RefreshToken, client)
	if err != nil {
		return err
	}

	//Generate tokens and Login the user
	accessToken, refreshToken, err := oauth.GenerateTokens(sessionRefreshToken.Client, sessionRefreshToken.User, sessionRefreshToken.Scope)
	if err != nil {
		return err
	}

	user.AccessToken = accessToken.Token
	user.RefreshToken = refreshToken.Token

	return nil
}
