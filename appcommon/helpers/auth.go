package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/TarsCloud/TarsGo/tars/util/current"
	"github.com/floppyisadog/appcommon/codes"
	"github.com/floppyisadog/appcommon/consts"
	"github.com/floppyisadog/appcommon/utils/crypto"
)

var (
	shortSession = time.Duration(12 * time.Hour)
	longSession  = time.Duration(30 * 24 * time.Hour)
)

// authorization
func GetAuth(ctx context.Context) (map[string]string, string, error) {
	metadata, ok := current.GetRequestContext(ctx)
	if !ok {
		return nil, "", codes.ErrRequestContextMissing
	}

	authz := metadata[consts.AuthorizationMetadata]
	if len(authz) == 0 {
		return nil, "", codes.ErrRequestAuthrizationMetaMissing
	}

	return metadata, authz, nil
}

// authentication
// LoginUser sets a cookie to log in a user
func LoginUser(uuid string, support, rememberMe bool, signingSecret, domain string, res http.ResponseWriter) {
	var dur time.Duration
	var maxAge int

	if rememberMe {
		// "Remember me"
		dur = longSession
		maxAge = 0
	} else {
		dur = shortSession
		maxAge = int(dur.Seconds())
	}
	token, err := crypto.SessionToken(uuid, signingSecret, support, dur)
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:   consts.CookieName,
		Value:  token,
		Path:   "/",
		Domain: "." + domain,
		MaxAge: maxAge,
	}
	http.SetCookie(res, cookie)
}

func GetSession(req *http.Request, signingSecret string) (uuid string, support bool, err error) {
	cookie, err := req.Cookie(consts.CookieName)
	if err != nil {
		return
	}
	uuid, support, err = crypto.RetrieveSessionInformation(cookie.Value, signingSecret)
	return
}

// Logout forces an immediate logout of the current user.
// For internal applications - typically do an HTTP redirect
// to the www service `/logout/` route to trigger a logout.
func Logout(res http.ResponseWriter, domain string) {
	cookie := &http.Cookie{
		Name:   consts.CookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		Domain: "." + domain,
	}
	http.SetCookie(res, cookie)
}
