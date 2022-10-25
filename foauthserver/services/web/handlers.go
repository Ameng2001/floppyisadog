package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/floppyisadog/foauthserver/common"
	"github.com/floppyisadog/foauthserver/models"
	"github.com/floppyisadog/foauthserver/services/oauth"
	"github.com/floppyisadog/foauthserver/util"
	"github.com/floppyisadog/foauthserver/util/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func registerFormHandler(c *gin.Context) {
	errMsg := util.OneFlash(c)
	util.RenderTemplate(c, "register.html", gin.H{
		"error":       errMsg,
		"queryString": util.GetQueryString(c.Request.URL.Query()),
	})
	// c.HTML(http.StatusOK, "register.html", gin.H{
	// 	"error":       errMsg,
	// 	"queryString": util.GetQueryString(c.Request.URL.Query()),
	// })
}

func registerHandler(c *gin.Context) {
	//check the user email hasnot been registered already
	if oauth.UserExists(c.PostForm("email")) {
		util.FlashMessage(c, "Email Taken")
		c.Redirect(http.StatusFound, c.Request.RequestURI)
		return
	}

	//create a user
	_, err := oauth.CreateUser(common.USER, c.PostForm("email"), c.PostForm("password"))
	if err != nil {
		util.FlashMessage(c, err.Error())
		c.Redirect(http.StatusFound, c.Request.RequestURI)
		return
	}

	//redirect to loginpage
	redirectURI := fmt.Sprintf("%s%s", "login", util.GetQueryString(c.Request.URL.Query()))
	c.Redirect(http.StatusFound, redirectURI)
}

func loginFormHandler(c *gin.Context) {
	errMsg := util.OneFlash(c)
	util.RenderTemplate(c, "login.html", gin.H{
		"error":       errMsg,
		"queryString": util.GetQueryString(c.Request.URL.Query()),
	})
	// c.HTML(http.StatusOK, "login.html", gin.H{
	// 	"error":       errMsg,
	// 	"queryString": util.GetQueryString(c.Request.URL.Query()),
	// })
}

func loginHandler(c *gin.Context) {
	session := sessions.Default(c)

	//Get the client
	client, err := oauth.FindClientByClientKey(c.Query("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("login", err))
		return
	}

	//Authenticate the user
	user, err := oauth.AuthUser(c.PostForm("email"), c.PostForm("password"))
	if err != nil {
		util.FlashMessage(c, err.Error())
		c.Redirect(http.StatusFound, c.Request.RequestURI)
		return
	}

	//Get the scope string
	scope, err := oauth.GetScope(c.Query("scope"))
	if err != nil {
		util.FlashMessage(c, err.Error())
		c.Redirect(http.StatusFound, c.Request.RequestURI)
		return
	}

	//Login the user
	accessToken, refreshToken, err := oauth.GenerateTokens(client, user, scope)
	if err != nil {
		util.FlashMessage(c, err.Error())
		c.Redirect(http.StatusFound, c.Request.RequestURI)
		return
	}

	// store the user session in a cookie
	userSession := &UserSession{
		ClientID:     client.ClientKey,
		Username:     user.Username,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}
	session.Set(USER_SESSION_KEY, userSession)
	session.Save()

	// redirect to a query string
	loginRedirectURI := c.Query("login_redirect_uri")
	if loginRedirectURI == "" {
		loginRedirectURI = "admin"
	}
	c.Redirect(http.StatusFound, loginRedirectURI)
}

func authorizeFormHandler(c *gin.Context) {
	errMsg := util.OneFlash(c)
	query := c.Request.URL.Query()
	query.Set("login_redirect_uri", c.Request.URL.Path)

	client, _, responseType, _, err := authorizeCommon(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("authorize", err))
		return
	}

	util.RenderTemplate(c, "authorize.html", gin.H{
		"error":       errMsg,
		"clientID":    client.ClientKey,
		"queryString": util.GetQueryString(c.Request.URL.Query()),
		"token":       responseType == "token",
	})
	// c.HTML(http.StatusOK, "authorize.html", gin.H{
	// 	"error":       errMsg,
	// 	"clientID":    client.ClientKey,
	// 	"queryString": util.GetQueryString(c.Request.URL.Query()),
	// 	"token":       responseType == "token",
	// })
}

func authorizeHandler(c *gin.Context) {
	client, user, responseType, redirectURI, err := authorizeCommon(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("authorize", err))
		return
	}

	//Get the state param
	state := c.Query("state")

	// Has the resource owner or authorization server denied the request?
	authorized := len(c.PostForm("allow")) > 0
	if !authorized {
		util.ErrorRedirect(c, redirectURI, "access_denied", state, responseType)
		return
	}

	//Check the requested scope
	scope, err := oauth.GetScope(c.PostForm("scope"))
	if err != nil {
		util.ErrorRedirect(c, redirectURI, "invalid_scope", state, responseType)
		return
	}

	query := redirectURI.Query()
	// When response_type == "code", we will grant an authorization code
	if responseType == "code" {
		// Create a new authorization code
		authcode, err := oauth.GrantAuthorizationCode(
			client,
			user,
			config.GetConfig().Oauth.AuthCodeLifetime,
			redirectURI.String(),
			scope,
		)
		if err != nil {
			util.ErrorRedirect(c, redirectURI, "server_error", state, responseType)
			return
		}

		// Set query string params for the redirection URL
		query.Set("code", authcode.Code)
		// Add state if present
		if state != "" {
			query.Set("state", state)
		}

		util.RedirectWithQueryString(redirectURI.String(), query, c)

		return
	}
	// When response_type == "token", we will directly grant an access token
	if responseType == "token" {
		lifeTime, err := strconv.Atoi(c.PostForm("lifetime"))
		if err != nil {
			util.ErrorRedirect(c, redirectURI, "server_error", state, responseType)
			return
		}

		// Grant access token directly
		accessToken, err := oauth.GrantAccessToken(client, user, lifeTime, scope)
		if err != nil {
			util.ErrorRedirect(c, redirectURI, "server_error", state, responseType)
			return
		}

		// Set query string params for the redirection URL
		query.Set("access_token", accessToken.Token)
		query.Set("expires_in", fmt.Sprintf("%d", lifeTime))
		query.Set("token_type", common.TOKEN_TYPE_BEARER)
		query.Set("scope", scope)
		// Add state param if present (recommended)
		if state != "" {
			query.Set("state", state)
		}

		util.RedirectWithQueryString(redirectURI.String(), query, c)
	}
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	sessionUser := session.Get(USER_SESSION_KEY).(*UserSession)

	//delete access and refresh token
	oauth.ClearUserTokens(sessionUser.AccessToken, sessionUser.RefreshToken)

	//clear session
	session.Delete(USER_SESSION_KEY)
	session.Save()

	util.RedirectWithQueryString("/login", c.Request.URL.Query(), c)
}

func authorizeCommon(c *gin.Context) (*models.OauthClient, *models.OauthUser, string, *url.URL, error) {
	session := sessions.Default(c)

	//Get the client
	client, err := oauth.FindClientByClientKey(c.Query("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("authorize", err))
		return nil, nil, "", nil, err
	}

	sessionUser := session.Get(USER_SESSION_KEY).(*UserSession)
	user, err := oauth.FindUserByUserName(sessionUser.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("authorize", err))
		return nil, nil, "", nil, err
	}

	//check the responsetype is either code or token
	responseType := c.Query("response_type")
	if responseType != "code" && responseType != "token" {
		return nil, nil, "", nil, common.ErrIncorrectResponseType
	}

	// fallback to the client redirect URI if not in query string
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = client.RedirectURL
	}
	parsedRedirectURI, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		return nil, nil, "", nil, err
	}

	return client, user, responseType, parsedRedirectURI, nil
}
