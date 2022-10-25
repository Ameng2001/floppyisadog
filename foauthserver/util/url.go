package util

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Redirects to a new path while keeping current request's query string
func RedirectWithQueryString(to string, query url.Values, c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("%s%s", to, GetQueryString(query)))
}

// Redirects to a new path with the query string moved to the URL fragment
func RedirectWithFragment(to string, query url.Values, c *gin.Context) {
	c.Redirect(http.StatusFound, fmt.Sprintf("%s#%s", to, query.Encode()))
}

// Helper function to handle redirecting failed or declined authorization
func ErrorRedirect(c *gin.Context, redirectURI *url.URL, err, state, responseType string) {
	query := redirectURI.Query()
	query.Set("error", err)
	if state != "" {
		query.Set("state", state)
	}
	if responseType == "code" {
		RedirectWithQueryString(redirectURI.String(), query, c)
	}
	if responseType == "token" {
		RedirectWithFragment(redirectURI.String(), query, c)
	}
}

// Helper function to save flash message
func FlashMessage(c *gin.Context, message string) {
	session := sessions.Default(c)
	session.AddFlash(message)
	if err := session.Save(); err != nil {
		log.Printf("error in save flash message : %s", err)
	}
}

// Helper function to retrive one flash message
func OneFlash(c *gin.Context) string {
	session := sessions.Default(c)
	flashes := session.Flashes()

	if len(flashes) != 0 {
		if err := session.Save(); err != nil {
			log.Printf("error in saving the session %s", err)
		}
		return flashes[0].(string)
	}

	return ""
}

// Helper function to rtirve all the flash messages
func AllFlashes(c *gin.Context) []interface{} {
	session := sessions.Default(c)
	flashes := session.Flashes()

	if len(flashes) != 0 {
		if err := session.Save(); err != nil {
			log.Printf("error in saving the session %s", err)
		}
	}

	return flashes
}
