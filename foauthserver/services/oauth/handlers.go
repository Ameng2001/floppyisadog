package oauth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/floppyisadog/foauthserver/common"
	"github.com/floppyisadog/foauthserver/models"
	"github.com/floppyisadog/foauthserver/util"
	"github.com/floppyisadog/foauthserver/util/config"
	"github.com/gin-gonic/gin"
)

// Map of grant types against handler functions
var (
	grantTypes = map[string]func(c *gin.Context, client *models.OauthClient) (*AccessTokenResponse, error){
		"authorization_code": authorizationCodeGrant,
		"password":           passwordGrant,
		"client_credentials": clientCredentialsGrant,
		"refresh_token":      refreshTokenGrant,
	}
)

func tokensHandler(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("oauth", err))
		return
	}

	grantHandler, ok := grantTypes[c.Request.Form.Get("grant_type")]
	if !ok {
		c.JSON(http.StatusBadRequest, common.NewError("oauth", errors.New("invalid grant type")))
		return
	}

	// client auth
	client, err := basicAuthClient(c)
	if err != nil {
		c.Header("WWW-Authenticate", fmt.Sprintf("Bearer realm=%s", "floppyisadog-oauth2-server"))
		c.JSON(http.StatusUnauthorized, common.NewError("oauth", err))
		return
	}

	// grant processing
	resp, err := grantHandler(c, client)
	if err != nil {
		c.JSON(common.GetErrStatusCode(err), common.NewError("oauth", err))
		return
	}

	c.JSON(http.StatusOK, *resp)
}

func introspectHandler(c *gin.Context) {
	// client auth
	client, err := basicAuthClient(c)
	if err != nil {
		c.Header("WWW-Authenticate", fmt.Sprintf("Bearer realm=%s", "floppyisadog-oauth2-server"))
		c.JSON(http.StatusUnauthorized, common.NewError("oauth", err))
		return
	}

	// introspect the token
	resp, err := introspectToken(c, client)
	if err != nil {
		c.JSON(common.GetErrStatusCode(err), common.NewError("oauth", err))
		return
	}

	c.JSON(http.StatusOK, *resp)
}

func authorizationCodeGrant(c *gin.Context, client *models.OauthClient) (*AccessTokenResponse, error) {
	authCode, err := getValidAuthorizationCode(
		c.Request.FormValue("code"),
		c.Request.FormValue("redirect_uri"),
		client,
	)
	if err != nil {
		return nil, err
	}

	//Login the user
	accessToken, refreshToken, err := GenerateTokens(authCode.Client, authCode.User, authCode.Scope)
	if err != nil {
		return nil, err
	}

	//Release the authcode
	deleteAuthorizationCode(authCode)

	//Create response
	serializer := AccessTokenSerializer{
		accessToken,
		refreshToken,
		config.GetConfig().Oauth.AccessTokenLifetime,
		common.TOKEN_TYPE_BEARER,
	}

	return serializer.Response()
}

func passwordGrant(c *gin.Context, client *models.OauthClient) (*AccessTokenResponse, error) {
	// Get the scope form form param
	scope, err := GetScope(c.Request.FormValue("scope"))
	if err != nil {
		return nil, err
	}

	// Authenticate the user
	user, err := AuthUser(c.Request.FormValue("username"), c.Request.FormValue("password"))
	if err != nil {
		return nil, common.ErrInvalidUsernameOrPassword
	}

	// login the user
	accessToken, refreshToken, err := GenerateTokens(client, user, scope)
	if err != nil {
		return nil, err
	}

	//Create response
	serializer := AccessTokenSerializer{
		accessToken,
		refreshToken,
		config.GetConfig().Oauth.AccessTokenLifetime,
		common.TOKEN_TYPE_BEARER,
	}

	return serializer.Response()
}

func clientCredentialsGrant(c *gin.Context, client *models.OauthClient) (*AccessTokenResponse, error) {
	// Get the scope form form param
	scope, err := GetScope(c.Request.FormValue("scope"))
	if err != nil {
		return nil, err
	}

	accessToken, err := GrantAccessToken(client, nil, config.GetConfig().Oauth.AccessTokenLifetime, scope)
	if err != nil {
		return nil, err
	}

	//Create response
	serializer := AccessTokenSerializer{
		accessToken,
		nil,
		config.GetConfig().Oauth.AccessTokenLifetime,
		common.TOKEN_TYPE_BEARER,
	}

	return serializer.Response()
}

func refreshTokenGrant(c *gin.Context, client *models.OauthClient) (*AccessTokenResponse, error) {
	//Fetch the refresh token
	refreshToken, err := ValidateRefreshToken(c.Request.FormValue("refresh_token"), client)
	if err != nil {
		return nil, err
	}

	//Get the scope
	scope, err := getRefreshTokenScope(refreshToken, c.Request.FormValue("scope"))
	if err != nil {
		return nil, err
	}

	//login the user
	accessToken, refreshToken, err := GenerateTokens(
		refreshToken.Client,
		refreshToken.User,
		scope,
	)
	if err != nil {
		return nil, err
	}

	//Create response
	serializer := AccessTokenSerializer{
		accessToken,
		refreshToken,
		config.GetConfig().Oauth.AccessTokenLifetime,
		common.TOKEN_TYPE_BEARER,
	}

	return serializer.Response()
}

func basicAuthClient(c *gin.Context) (*models.OauthClient, error) {
	// Get client credentials from basic auth
	name, pwd, ok := c.Request.BasicAuth()
	if !ok {
		return nil, common.ErrInvalidClientIDOrSecret
	}

	// Authenticate the client
	client, err := authClient(name, pwd)
	if err != nil {
		return nil, common.ErrInvalidClientIDOrSecret
	}

	return client, nil
}

func AuthUser(uname, upwd string) (*models.OauthUser, error) {
	//Fetch the user
	user, err := FindUserByUserName(uname)
	if err != nil {
		return nil, err
	}

	//Verify the password
	if !user.Password.Valid {
		return nil, common.ErrUserPasswordNotSet
	}
	if util.VerifyPassword(user.Password.String, upwd) != nil {
		return nil, common.ErrInvalidUserPassword
	}

	return user, nil
}

func GenerateTokens(client *models.OauthClient, user *models.OauthUser, scope string) (*models.OauthAccessToken, *models.OauthRefreshToken, error) {
	// Return error if user's role is not allowed to use this service
	// if !isRoleAllowed(user.RoleID.String) {
	// 	// For security reasons, return a general error message
	// 	return nil, nil, common.ErrInvalidUsernameOrPassword
	// }

	//Create a new access token
	accessToken, err := GrantAccessToken(
		client,
		user,
		config.GetConfig().Oauth.AccessTokenLifetime,
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := grantRefreshToken(
		client,
		user,
		config.GetConfig().Oauth.RefreshTokenLifetime,
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func introspectToken(c *gin.Context, client *models.OauthClient) (*IntrospectResponse, error) {
	// Parse the form so r.Form becomes available
	if err := c.Request.ParseForm(); err != nil {
		return nil, err
	}

	// Get token from the query
	token := c.Request.FormValue("token")
	if token == "" {
		return nil, common.ErrTokenMissing
	}

	// Get token type hint from the query
	tokenTypeHint := c.Request.FormValue("token_type_hint")
	if tokenTypeHint == "" {
		tokenTypeHint = common.ACCESSTOKEN_HINT
	}

	var accessToken *models.OauthAccessToken
	var refreshToken *models.OauthRefreshToken
	var err error

	switch tokenTypeHint {
	case common.ACCESSTOKEN_HINT:
		accessToken, err = ValidateAccessToken(token)
	case common.REFRESHTOKEN_HINT:
		refreshToken, err = ValidateRefreshToken(token, client)
	default:
		err = common.ErrTokenHintInvalid
	}

	if err != nil {
		return nil, err
	}

	serializer := IntrospectSerializer{
		tokenTypeHint,
		accessToken,
		refreshToken,
	}

	return serializer.Response()
}
