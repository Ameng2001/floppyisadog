package common

import (
	"errors"
	"net/http"
)

var (
	ErrAuthorizationCodeNotFound     = errors.New("authorization code not found")
	ErrAuthorizationCodeExpired      = errors.New("authorization code expired")
	ErrInvalidRedirectURI            = errors.New("invalid redirect URI")
	ErrInvalidScope                  = errors.New("invalid scope")
	ErrInvalidUsernameOrPassword     = errors.New("invalid username or password")
	ErrRefreshTokenNotFound          = errors.New("refresh token not found")
	ErrRefreshTokenExpired           = errors.New("refresh token expired")
	ErrRequestedScopeCannotBeGreater = errors.New("request scope error")
	ErrTokenMissing                  = errors.New("lack of token")
	ErrTokenHintInvalid              = errors.New("invalid token hint")
	ErrAccessTokenNotFound           = errors.New("accesstoken not found")
	ErrAccessTokenExpired            = errors.New("accesstoken expired")
	ErrInvalidClientIDOrSecret       = errors.New("invalid client ID or secret")
	ErrClientNotFound                = errors.New("client not found")
	ErrInvalidClientSecret           = errors.New("invalid client secret")
	ErrClientIDTaken                 = errors.New("client ID taken")
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserPasswordNotSet            = errors.New("user password not set")
	ErrInvalidUserPassword           = errors.New("invalid user password")
	ErrRoleNotFound                  = errors.New("role not found")
	ErrPasswordTooShort              = errors.New("password must be at least 6 characters long")
	ErrUsernameTaken                 = errors.New("username taken")
	ErrIncorrectResponseType         = errors.New("response type not one of token or code")
)

// define inner errors to standard http error codes
var (
	errStatusCodeMap = map[error]int{
		ErrAuthorizationCodeNotFound:     http.StatusNotFound,
		ErrAuthorizationCodeExpired:      http.StatusBadRequest,
		ErrInvalidRedirectURI:            http.StatusBadRequest,
		ErrInvalidScope:                  http.StatusBadRequest,
		ErrInvalidUsernameOrPassword:     http.StatusBadRequest,
		ErrRefreshTokenNotFound:          http.StatusNotFound,
		ErrRefreshTokenExpired:           http.StatusBadRequest,
		ErrRequestedScopeCannotBeGreater: http.StatusBadRequest,
		ErrTokenMissing:                  http.StatusNotFound,
		ErrTokenHintInvalid:              http.StatusBadRequest,
		ErrAccessTokenNotFound:           http.StatusNotFound,
	}
)

// My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"error"`
}

// Warp the error info in a object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

func GetErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	return http.StatusInternalServerError
}
