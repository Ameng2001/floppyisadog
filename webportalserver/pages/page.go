package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

const (
	// For SEO / web crawlers
	defaultDescription = "Staffjoy is an application that helps businesses create and share schedules with hourly workers."
)

type Page struct {
	Title        string // Used in <title>
	Description  string // SEO matters
	TemplateName string // e.g. home.html
	CSSId        string // e.g. 'careers'
	Version      string // e.g. master-1, for cachebusting
	CsrfField    template.HTML
}

type ActivatePage struct {
	Page
	ErrorMessage string
	Email        string
	Name         string
	Phonenumber  string
}

type BreaktimeEpisodePage struct {
	Page
	// Message stuff
	Name              string
	SoundcloudTrackID string
	Body              template.HTML
	CoverPhoto        string
	Date              string
}

type BreaktimeListPage struct {
	Page
	// Message stuff
	Episodes map[string]BreaktimeEpisodePage
}

type ConfirmResetPage struct {
	Page
	ErrorMessage string
}

type LoginPage struct {
	Page
	Denied        bool
	PreviousEmail string
	ReturnTo      string
}

type ResetPage struct {
	Page
	Denied          bool
	RecaptchaPublic string
}

func (p *Page) Handler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if p.Description == "" {
		p.Description = defaultDescription
	}

	p.CsrfField = middleware.TemplateField(c)

	err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p)

	if err != nil {
		fmt.Printf("Unable to render page %s - %s", p.Title, err)
	}
}
