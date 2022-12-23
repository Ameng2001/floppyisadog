package services

import (
	"fmt"
	"html/template"

	"github.com/floppyisadog/appcommon/errorpages"
	"github.com/floppyisadog/webportalserver/managers/assetsmgr"
	"github.com/floppyisadog/webportalserver/managers/configmgr"
	"github.com/floppyisadog/webportalserver/managers/logmgr"
	"github.com/floppyisadog/webportalserver/middleware"
	"github.com/gin-gonic/gin"
)

func breaktimeListHandler(c *gin.Context) {
	p := configmgr.GetPages().BreaktimeListPage
	p.CsrfField = middleware.TemplateField(c)
	p.Episodes = configmgr.GetPages().EpisodesPages

	if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, p.TemplateName, p); err != nil {
		panic(err)
	}
}

func breaktimeEpisodeHandler(c *gin.Context) {
	slug := c.Param("slug")

	episode, ok := configmgr.GetPages().EpisodesPages[slug]
	if !ok {
		errorpages.NotFound(c.Writer)
		return
	}

	body, ok := assetsmgr.GetBeakTime()[slug]
	if !ok {
		logmgr.RERROR("cannot find episode body for slug %s", slug)
		return
	}

	episode.Body = template.HTML(body)
	episode.CoverPhoto = fmt.Sprintf("/assets/breaktime-cover/%s.jpg", slug)

	if err := assetsmgr.GetTemplate().ExecuteTemplate(c.Writer, episode.TemplateName, episode); err != nil {
		panic(err)
	}
}
