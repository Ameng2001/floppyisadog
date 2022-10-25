package models

import (
	"database/sql"
	"time"

	"github.com/RichardKnop/uuid"
	"github.com/floppyisadog/foauthserver/util"
)

// oauth_clients table orm
type OauthClient struct {
	MyGormModel
	ClientKey    string `sql:"type:varchar(254);unique;not null"`
	ClientSecret string `sql:"type:varchar(60);not null"`
	RedirectURL  string `sql:"type:varchar(200)"`
}

func (c *OauthClient) TableName() string {
	return "oauth_clients"
}

// oauth_scopes table orm
type OauthScope struct {
	MyGormModel
	Scope       string `sql:"type:varchar(200);unique;not null"`
	Description sql.NullString
	IsDefault   bool `sql:"default:false"`
}

func (c *OauthScope) TableName() string {
	return "oauth_scopes"
}

// oauth_roles table orm
type OauthRole struct {
	TimestampModel
	ID   string `gorm:"primary_key" sql:"type:varchar(20)"`
	Name string `sql:"type:varchar(50);unique;not null"`
}

func (c *OauthRole) TableName() string {
	return "oauth_roles"
}

// oauth_users table orm
type OauthUser struct {
	MyGormModel
	RoleID   sql.NullString `sql:"type:varchar(20);index;not null"`
	Role     *OauthRole
	Username string         `sql:"type:varchar(254);unique;not null"`
	Password sql.NullString `sql:"type:varchar(60)"`
}

func (c *OauthUser) TableName() string {
	return "oauth_users"
}

// oauth_refresh_tokens table orm
type OauthRefreshToken struct {
	MyGormModel
	ClientID  sql.NullString `sql:"index;not null"`
	UserID    sql.NullString `sql:"index"`
	Client    *OauthClient
	User      *OauthUser
	Token     string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time `sql:"not null;DEFAULT:current_timestamp"`
	Scope     string    `sql:"type:varchar(200);not null"`
}

func (rt *OauthRefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

// oauth_access_tokens table orm
type OauthAccessToken struct {
	MyGormModel
	ClientID  sql.NullString `sql:"index;not null"`
	UserID    sql.NullString `sql:"index"`
	Client    *OauthClient
	User      *OauthUser
	Token     string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time `sql:"not null;DEFAULT:current_timestamp"`
	Scope     string    `sql:"type:varchar(200);not null"`
}

func (at *OauthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

// oauth_authorization_codes table orm
type OauthAuthorizationCode struct {
	MyGormModel
	ClientID    sql.NullString `sql:"index;not null"`
	UserID      sql.NullString `sql:"index;not null"`
	Client      *OauthClient
	User        *OauthUser
	Code        string         `sql:"type:varchar(40);unique;not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
	ExpiresAt   time.Time      `sql:"not null;DEFAULT:current_timestamp"`
	Scope       string         `sql:"type:varchar(200);not null"`
}

func (ac *OauthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

// the model factory methods
func NewOauthRefreshToken(client *OauthClient, user *OauthUser, expiresIn int, scope string) *OauthRefreshToken {
	refreshToken := &OauthRefreshToken{
		MyGormModel: MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		ClientID:  util.StringOrNull(string(client.ID)),
		Token:     uuid.New(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}

	if user != nil {
		refreshToken.UserID = util.StringOrNull(string(user.ID))
	}
	return refreshToken
}

func NewOauthAccessToken(client *OauthClient, user *OauthUser, expiresIn int, scope string) *OauthAccessToken {
	accessToken := &OauthAccessToken{
		MyGormModel: MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		ClientID:  util.StringOrNull(string(client.ID)),
		Token:     uuid.New(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		accessToken.UserID = util.StringOrNull(string(user.ID))
	}
	return accessToken
}

func NewOauthAuthorizationCode(client *OauthClient, user *OauthUser, expiresIn int, redirectURI, scope string) *OauthAuthorizationCode {
	return &OauthAuthorizationCode{
		MyGormModel: MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		ClientID:    util.StringOrNull(string(client.ID)),
		UserID:      util.StringOrNull(string(user.ID)),
		Code:        uuid.New(),
		ExpiresAt:   time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		RedirectURI: util.StringOrNull(redirectURI),
		Scope:       scope,
	}
}

//
