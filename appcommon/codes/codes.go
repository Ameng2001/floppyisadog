package codes

import (
	"errors"
	"net/http"
)

const (
	OK                 int32 = 100
	Canceled           int32 = 1
	Unknown            int32 = 2
	InvalidArgument    int32 = 3
	DeadlineExceeded   int32 = 4
	NotFound           int32 = 5
	AlreadyExists      int32 = 6
	PermissionDenied   int32 = 7
	ResourceExhausted  int32 = 8
	FailedPrecondition int32 = 9
	Aborted            int32 = 10
	OutOfRange         int32 = 11
	Unimplemented      int32 = 12
	Internal           int32 = 13
	Unavailable        int32 = 14
	DataLoss           int32 = 15
	Unauthenticated    int32 = 16
)

var (
	//Errors for Oauth
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

	//errors for authhelper
	ErrRequestContextMissing          = errors.New("request context missing")
	ErrRequestAuthrizationMetaMissing = errors.New("authrization meta missing")
	ErrAuthorizedFailed               = errors.New("failed to authorize")
	ErrPermissionDenied               = errors.New("permission denied")
	ErrInvalidArgument                = errors.New("invalid arguments")
	ErrAccountNotFound                = errors.New("account not found")
	ErrAccountAlreadyExists           = errors.New("account already exists")
	ErrGenerateUUID                   = errors.New("can not generate uuid")
	ErrCreateAccount                  = errors.New("create account error")
	ErrSendActiveEmail                = errors.New("send active email error")

	//errors for account
	ErrUpdateAccountError = errors.New("update account error")
	ErrInternal           = errors.New("internal error")
)

// define inner errors to standard http error codes
var (
	errStatusCodeMap = map[error]int{
		// Inner errors for OAuth
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
