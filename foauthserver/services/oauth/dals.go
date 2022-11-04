package oauth

import (
	"sort"
	"strings"
	"time"

	"github.com/RichardKnop/uuid"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/utils"
	"github.com/floppyisadog/appcommon/utils/database"
	"github.com/floppyisadog/foauthserver/managers/configmgr"
	"github.com/floppyisadog/foauthserver/models"
	"github.com/jinzhu/gorm"
)

// Client
func authClient(clientID, secret string) (*models.OauthClient, error) {
	// Fetch the client
	client, err := FindClientByClientKey(clientID)
	if err != nil {
		return nil, codes.ErrClientNotFound
	}

	// Verify the secret
	if utils.VerifyPassword(client.ClientSecret, secret) != nil {
		return nil, codes.ErrInvalidClientSecret
	}

	return client, nil
}

func FindClientByClientKey(clientID string) (*models.OauthClient, error) {
	client := new(models.OauthClient)
	notFound := database.GetDB().Where("client_key = LOWER(?)", clientID).First(client).RecordNotFound()

	if notFound {
		return nil, codes.ErrClientNotFound
	}

	return client, nil
}

func FindClientByID(id string) (*models.OauthClient, error) {
	client := new(models.OauthClient)
	nofound := database.GetDB().Select("client_key").Where("id = ?", id).
		First(client).RecordNotFound()

	if nofound {
		return nil, codes.ErrUserNotFound
	}

	return client, nil
}

// User
func UserExists(username string) bool {
	_, err := FindUserByUserName(username)
	return err == nil
}

func CreateUser(roleid, username, password string) (*models.OauthUser, error) {
	//start a user without a password
	user := &models.OauthUser{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		RoleID:   utils.StringOrNull(roleid),
		Username: strings.ToLower(username),
		Password: utils.StringOrNull(""),
	}

	if password != "" {
		if len(password) < configmgr.GetConfig().MinPwdLength {
			return nil, codes.ErrPasswordTooShort
		}

		passwordHash, err := utils.HashPassword(password)
		if err != nil {
			return nil, err
		}

		user.Password = utils.StringOrNull(string(passwordHash))
	}

	//check the username is valid
	if UserExists(user.Username) {
		return nil, codes.ErrUsernameTaken
	}

	//create the user in the db
	if err := database.GetDB().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByUserName(usrname string) (*models.OauthUser, error) {
	user := new(models.OauthUser)
	nofound := database.GetDB().Where("username = LOWER(?)", usrname).
		First(user).RecordNotFound()

	if nofound {
		return nil, codes.ErrUserNotFound
	}

	return user, nil
}

func findUserByUserID(id string) (*models.OauthUser, error) {
	user := new(models.OauthUser)
	nofound := database.GetDB().Select("username").Where("id = ?", id).
		First(user).RecordNotFound()

	if nofound {
		return nil, codes.ErrUserNotFound
	}

	return user, nil
}

// Role
// func findRoleByID(id string) (*models.OauthRole, error) {
// 	role := new(models.OauthRole)
// 	if database.GetDB().Where("id = ?", id).First(role).RecordNotFound() {
// 		return nil, codes.ErrRoleNotFound
// 	}

// 	return role, nil
// }

// AuthorizationCode
func GrantAuthorizationCode(client *models.OauthClient, user *models.OauthUser, expiredIn int, redirectURI, scope string) (*models.OauthAuthorizationCode, error) {
	authcode := models.NewOauthAuthorizationCode(client, user, expiredIn, redirectURI, scope)
	if err := database.GetDB().Create(authcode).Error; err != nil {
		return nil, err
	}

	authcode.Client = client
	authcode.User = user

	return authcode, nil
}

func getValidAuthorizationCode(code, redirectURL string, client *models.OauthClient) (*models.OauthAuthorizationCode, error) {
	authCode := new(models.OauthAuthorizationCode)
	nofound := database.GetDB().
		Preload("Client").
		Preload("User").
		Where("client_id = ?", client.ID).
		Where("code = ?", code).First(authCode).RecordNotFound()
	if nofound {
		return nil, codes.ErrAuthorizationCodeNotFound
	}

	if redirectURL != authCode.RedirectURI.String {
		return nil, codes.ErrInvalidRedirectURI
	}

	if time.Now().After(authCode.ExpiresAt) {
		return nil, codes.ErrAuthorizationCodeExpired
	}

	return authCode, nil
}

func deleteAuthorizationCode(code *models.OauthAuthorizationCode) {
	database.GetDB().Unscoped().Delete(*code)
}

// Scope
func GetScope(requestedScope string) (string, error) {
	if requestedScope == "" {
		return getDefaultScope(), nil
	}

	if scopeExists(requestedScope) {
		return requestedScope, nil
	}

	return "", codes.ErrInvalidScope
}

func getDefaultScope() string {
	var scopes []string
	database.GetDB().Model(new(models.OauthScope)).
		Where("is_default = ?", true).
		Pluck("scope", &scopes)

	sort.Strings(scopes)

	return strings.Join(scopes, " ")
}

func scopeExists(requestScope string) bool {
	scopes := strings.Split(requestScope, " ")

	var count int
	database.GetDB().Model(new(models.OauthScope)).
		Where("scope in (?)", scopes).
		Count(&count)

	return count == len(scopes)
}

// AccessToken
func GrantAccessToken(client *models.OauthClient, user *models.OauthUser, expiresIn int, scope string) (*models.OauthAccessToken, error) {
	tx := database.GetDB().Begin()

	query := tx.Unscoped().Where("client_id = ?", client.ID)
	if user != nil && len([]rune(user.ID)) > 0 {
		query = query.Where("user_id = ?", user.ID)
	} else {
		query = query.Where("user_id IS NULL")
	}

	if err := query.Where("expires_at <= ?", time.Now()).Delete(new(models.OauthAccessToken)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	accessToken := models.NewOauthAccessToken(client, user, expiresIn, scope)
	if err := tx.Create(accessToken).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	accessToken.Client = client
	accessToken.User = user

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return accessToken, nil
}

// Authenticate
func ValidateAccessToken(token string) (*models.OauthAccessToken, error) {
	accessToken := new(models.OauthAccessToken)
	nofound := database.GetDB().Where("token = ?", token).
		First(accessToken).
		RecordNotFound()

	if nofound {
		return nil, codes.ErrAccessTokenNotFound
	}

	if time.Now().UTC().After(accessToken.ExpiresAt) {
		return nil, codes.ErrAccessTokenExpired
	}

	query := database.GetDB().Model(new(models.OauthRefreshToken)).
		Where("client_id = ?", accessToken.ClientID.String)

	if accessToken.UserID.Valid {
		query = query.Where("user_id = ?", accessToken.UserID.String)
	} else {
		query = query.Where("user_id IS NULL")
	}

	increasedExpiresAt := gorm.NowFunc().Add(
		time.Duration(configmgr.GetConfig().Oauth.RefreshTokenLifetime) * time.Second,
	)

	if err := query.UpdateColumn("expires_at", increasedExpiresAt).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// RefreshToken
func grantRefreshToken(client *models.OauthClient, user *models.OauthUser, expiresIn int, scope string) (*models.OauthRefreshToken, error) {
	refreshToken := new(models.OauthRefreshToken)
	query := database.GetDB().
		Preload("Client").
		Preload("User").
		Where("client_id = ?", client.ID)

	if user != nil && len([]rune(user.ID)) > 0 {
		query = query.Where("user_id = ?", user.ID)
	} else {
		query = query.Where("user_id IS NULL")
	}

	found := !query.First(refreshToken).RecordNotFound()
	var expired bool
	if found {
		expired = time.Now().UTC().After(refreshToken.ExpiresAt)
	}

	if expired {
		database.GetDB().Delete(refreshToken)
	}

	if expired || !found {
		refreshToken = models.NewOauthRefreshToken(client, user, expiresIn, scope)
		if err := database.GetDB().Create(refreshToken).Error; err != nil {
			return nil, err
		}
		refreshToken.Client = client
		refreshToken.User = user
	}

	return refreshToken, nil
}

func ValidateRefreshToken(token string, client *models.OauthClient) (*models.OauthRefreshToken, error) {
	refreshToken := new(models.OauthRefreshToken)
	notfound := database.GetDB().
		Preload("Client").
		Preload("User").
		Where("client_id = ?", client.ID).
		Where("token = ?", token).
		First(refreshToken).
		RecordNotFound()

	if notfound {
		return nil, codes.ErrRefreshTokenNotFound
	}

	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		return nil, codes.ErrRefreshTokenExpired
	}

	return refreshToken, nil
}

func getRefreshTokenScope(refreshToken *models.OauthRefreshToken, requestScope string) (string, error) {
	var (
		scope = refreshToken.Scope
		err   error
	)

	if requestScope != "" {
		scope, err = GetScope(requestScope)
		if err != nil {
			return "", err
		}
	}

	if !utils.SpaceDelimitedStringNotGreater(scope, refreshToken.Scope) {
		return "", codes.ErrRequestedScopeCannotBeGreater
	}

	return scope, nil
}

// Clear Tokens
func ClearUserTokens(accessToken, refreshToken string) {
	// Clear all refresh tokens with user_id and client_id
	refreshTokenObj := new(models.OauthRefreshToken)
	found := !database.GetDB().
		Preload("Client").
		Preload("User").
		Where("token = ?", refreshToken).
		First(refreshTokenObj).
		RecordNotFound()

	if found {
		database.GetDB().
			Unscoped().
			Where("client_id = ? AND user_id = ?", refreshTokenObj.ClientID, refreshTokenObj.UserID).
			Delete(models.OauthRefreshToken{})
	}

	// Clear all access tokens with user_id and client_id
	accessTokenObj := new(models.OauthAccessToken)
	found = !database.GetDB().
		Preload("Client").
		Preload("User").
		Where("token = ?", accessToken).
		First(accessTokenObj).
		RecordNotFound()

	if found {
		database.GetDB().
			Unscoped().
			Where("client_id = ? AND user_id = ?", accessTokenObj.ClientID, accessTokenObj.UserID).
			Delete(models.OauthAccessToken{})
	}

}
