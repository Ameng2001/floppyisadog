package oauth

import (
	"github.com/floppyisadog/appcommon/enums"
	"github.com/floppyisadog/foauthserver/models"
)

// AccessTokenResponse ...
type AccessTokenSerializer struct {
	AccessTokenModel  *models.OauthAccessToken
	RefreshTokenModel *models.OauthRefreshToken
	LifeTime          int
	TokenType         string
}

type AccessTokenResponse struct {
	UserID       string `json:"user_id,omitempty"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (s *AccessTokenSerializer) Response() (*AccessTokenResponse, error) {
	response := &AccessTokenResponse{
		AccessToken: s.AccessTokenModel.Token,
		ExpiresIn:   s.LifeTime,
		TokenType:   s.TokenType,
		Scope:       s.AccessTokenModel.Scope,
	}

	if s.AccessTokenModel.UserID.Valid {
		response.UserID = s.AccessTokenModel.UserID.String
	}
	if s.RefreshTokenModel != nil {
		response.RefreshToken = s.RefreshTokenModel.Token
	}

	return response, nil
}

// IntrospectResponse ...
type IntrospectSerializer struct {
	tokenType         string
	accessTokenModel  *models.OauthAccessToken
	refreshTokenModel *models.OauthRefreshToken
}

type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Username  string `json:"username,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	ExpiresAt int    `json:"exp,omitempty"`
}

func (s *IntrospectSerializer) Response() (*IntrospectResponse, error) {
	var introspectResponse = &IntrospectResponse{
		Active:    true,
		TokenType: enums.TOKEN_TYPE_BEARER,
	}

	if s.tokenType == enums.ACCESSTOKEN_HINT {
		introspectResponse.Scope = s.accessTokenModel.Scope
		introspectResponse.ExpiresAt = int(s.accessTokenModel.ExpiresAt.Unix())
		if s.accessTokenModel.ClientID.Valid {
			client, err := FindClientByID(s.accessTokenModel.ClientID.String)
			if err != nil {
				return nil, err
			}
			introspectResponse.ClientID = client.ClientKey
		}

		if s.accessTokenModel.UserID.Valid {
			user, err := findUserByUserID(s.accessTokenModel.UserID.String)
			if err != nil {
				return nil, err
			}
			introspectResponse.Username = user.Username
		}
	} else if s.tokenType == enums.REFRESHTOKEN_HINT {
		introspectResponse.Scope = s.refreshTokenModel.Scope
		introspectResponse.ExpiresAt = int(s.refreshTokenModel.ExpiresAt.Unix())
		if s.refreshTokenModel.ClientID.Valid {
			client, err := FindClientByID(s.refreshTokenModel.ClientID.String)
			if err != nil {
				return nil, err
			}
			introspectResponse.ClientID = client.ClientKey
		}

		if s.refreshTokenModel.UserID.Valid {
			user, err := findUserByUserID(s.refreshTokenModel.UserID.String)
			if err != nil {
				return nil, err
			}
			introspectResponse.Username = user.Username
		}
	}

	return introspectResponse, nil
}
