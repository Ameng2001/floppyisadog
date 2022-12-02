package middleware

import (
	"fmt"
	"html/template"

	"github.com/floppyisadog/appcommon/errorpages"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const (
	csrfRequestHeader string = "X-CSRF-TOKEN"
	csrfFiled         string = "csrf"
	cookieStoreSecret string = "floppy-cookie-secret"
	sessionStoreName  string = "floppy-sessionstore"
)

var customIgnoreMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

var customErrorFunc = func(c *gin.Context) {
	//fmt.Printf("failed CSRF - %s", csrf.FailureReason(req))
	fmt.Printf("failed CSRF")
	errorpages.Forbidden(c.Writer)
}

var customTokenGetter = func(c *gin.Context) string {
	r := c.Request

	// 1. Check the HTTP header first.
	issued := r.Header.Get(csrfRequestHeader)

	// 2. Fall back to the POST (form) value.
	if issued == "" {
		issued = r.PostFormValue(csrfFiled)
	}

	// 3. Finally, fall back to the multipart form (if set).
	if issued == "" && r.MultipartForm != nil {
		vals := r.MultipartForm.Value[csrfFiled]

		if len(vals) > 0 {
			issued = vals[0]
		}
	}

	return issued
}

func InitCSRF(r *gin.Engine, secret string) {
	store := cookie.NewStore([]byte(cookieStoreSecret))
	r.Use(sessions.Sessions(sessionStoreName, store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret:        secret,
		IgnoreMethods: customIgnoreMethods,
		ErrorFunc:     customErrorFunc,
		TokenGetter:   customTokenGetter,
	}))
}

func TemplateField(c *gin.Context) template.HTML {
	fragment := fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`,
		csrfFiled, csrf.GetToken(c))

	return template.HTML(fragment)
}
