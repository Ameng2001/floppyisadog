package configmgr

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars/util/conf"
	"github.com/floppyisadog/appcommon/utils/environment"
	"github.com/floppyisadog/webportalserver/pages"
)

type PageConfig struct {
	Version           string
	StaticPages       map[string]*pages.Page
	ConfirmPage       *pages.Page
	ResetConfirmPage  *pages.Page
	NewCompanyPage    *pages.Page
	ActivatePage      *pages.ActivatePage
	LoginPage         *pages.LoginPage
	ConfirmResetPage  *pages.ConfirmResetPage
	ResetPage         *pages.ResetPage
	BreaktimeListPage *pages.BreaktimeListPage
	EpisodesPages     map[string]*pages.BreaktimeEpisodePage
}

// Config stores all configuration options
type Config struct {
	IsDevelopment bool `default:"True"`
	Outerfactory  map[string]string
	Pages         PageConfig
}

// DefaultConfig ...
// Let's start with some sensible defaults
var (
	ConfigInstance *Config
	defaultConfig  = &Config{
		IsDevelopment: true,
	}
)

func InitConfig(configFile string) {
	if configFile != "" {
		ConfigInstance = &Config{}
		c, err := conf.NewConf(configFile)
		if err != nil {
			//log.RERROR("Parse server config fail", err)
			fmt.Printf("Parse server config fail, err:(%s)\n", err)
			ConfigInstance = nil
		}

		ConfigInstance.IsDevelopment = c.GetBoolWithDef("/floppyisadog/webportalserver/<IsDevelopment>", true)

		environment.InitEnvironment(c, "/floppyisadog/environment/")

		ConfigInstance.Outerfactory = c.GetMap("/floppyisadog/webportalserver/outerfactory/")
		fmt.Printf("Init Proxy config: (%v)\n", ConfigInstance.Outerfactory)

		//parse page config
		ConfigInstance.Pages.Version = c.GetString("/floppyisadog/webportalserver/pages/<Version>")
		fmt.Printf("ConfigInstance.Pages.Version = (%s)", ConfigInstance.Pages.Version)
		staticPages := c.GetDomain("/floppyisadog/webportalserver/pages/staticpages")
		ConfigInstance.Pages.StaticPages = make(map[string]*pages.Page)
		for _, page := range staticPages {
			pageInfo := c.GetMap("/floppyisadog/webportalserver/pages/staticpages/" + page)
			ConfigInstance.Pages.StaticPages[pageInfo["PATH"]] = &pages.Page{
				Title:        pageInfo["Title"],
				Description:  pageInfo["Description"],
				TemplateName: pageInfo["TemplateName"],
				CSSId:        pageInfo["CSSId"],
				Version:      ConfigInstance.Pages.Version,
			}
			fmt.Printf("Init static page (%s) = (%v)\n", page, pageInfo)
		}
		fmt.Printf("Init static pages: (%v)\n", ConfigInstance.Pages.StaticPages)

		pageInfo := c.GetMap("/floppyisadog/webportalserver/pages/confirm")
		ConfigInstance.Pages.ConfirmPage = &pages.Page{
			Title:        pageInfo["Title"],
			Description:  pageInfo["Description"],
			TemplateName: pageInfo["TemplateName"],
			CSSId:        pageInfo["CSSId"],
			Version:      ConfigInstance.Pages.Version,
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/resetconfirm")
		ConfigInstance.Pages.ResetConfirmPage = &pages.Page{
			Title:        pageInfo["Title"],
			Description:  pageInfo["Description"],
			TemplateName: pageInfo["TemplateName"],
			CSSId:        pageInfo["CSSId"],
			Version:      ConfigInstance.Pages.Version,
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/newcompany")
		ConfigInstance.Pages.NewCompanyPage = &pages.Page{
			Title:        pageInfo["Title"],
			Description:  pageInfo["Description"],
			TemplateName: pageInfo["TemplateName"],
			CSSId:        pageInfo["CSSId"],
			Version:      ConfigInstance.Pages.Version,
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/login")
		ConfigInstance.Pages.LoginPage = &pages.LoginPage{
			Page: pages.Page{
				Title:        pageInfo["Title"],
				Description:  pageInfo["Description"],
				TemplateName: pageInfo["TemplateName"],
				CSSId:        pageInfo["CSSId"],
				Version:      ConfigInstance.Pages.Version,
			},
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/activate")
		ConfigInstance.Pages.ActivatePage = &pages.ActivatePage{
			Page: pages.Page{
				Title:        pageInfo["Title"],
				Description:  pageInfo["Description"],
				TemplateName: pageInfo["TemplateName"],
				CSSId:        pageInfo["CSSId"],
				Version:      ConfigInstance.Pages.Version,
			},
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/confirmreset")
		ConfigInstance.Pages.ConfirmResetPage = &pages.ConfirmResetPage{
			Page: pages.Page{
				Title:        pageInfo["Title"],
				Description:  pageInfo["Description"],
				TemplateName: pageInfo["TemplateName"],
				CSSId:        pageInfo["CSSId"],
				Version:      ConfigInstance.Pages.Version,
			},
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/password-reset")
		ConfigInstance.Pages.ResetPage = &pages.ResetPage{
			Page: pages.Page{
				Title:        pageInfo["Title"],
				Description:  pageInfo["Description"],
				TemplateName: pageInfo["TemplateName"],
				CSSId:        pageInfo["CSSId"],
				Version:      ConfigInstance.Pages.Version,
			},
		}

		pageInfo = c.GetMap("/floppyisadog/webportalserver/pages/breaktime-list")
		ConfigInstance.Pages.BreaktimeListPage = &pages.BreaktimeListPage{
			Page: pages.Page{
				Title:        pageInfo["Title"],
				Description:  pageInfo["Description"],
				TemplateName: pageInfo["TemplateName"],
				CSSId:        pageInfo["CSSId"],
				Version:      ConfigInstance.Pages.Version,
			},
		}

		episodePages := c.GetDomain("/floppyisadog/webportalserver/pages/breaktime-episodes")
		ConfigInstance.Pages.EpisodesPages = make(map[string]*pages.BreaktimeEpisodePage)
		for _, page := range episodePages {
			pageInfo := c.GetMap("/floppyisadog/webportalserver/pages/breaktime-episodes/" + page)
			ConfigInstance.Pages.EpisodesPages[page] = &pages.BreaktimeEpisodePage{
				Page: pages.Page{
					Title:        pageInfo["Title"],
					Description:  pageInfo["Description"],
					TemplateName: pageInfo["TemplateName"],
					CSSId:        pageInfo["CSSId"],
					Version:      ConfigInstance.Pages.Version,
				},
				SoundcloudTrackID: pageInfo["SoundcloudTrackID"],
				Date:              pageInfo["Date"],
			}
		}

	}
}

func GetConfig() *Config {
	if ConfigInstance == nil {
		return defaultConfig
	}

	return ConfigInstance
}

func GetPages() *PageConfig {
	return &ConfigInstance.Pages
}
