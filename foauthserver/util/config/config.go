package config

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars/util/conf"
	"github.com/floppyisadog/foauthserver/util"
)

// DatabaseConfig stores database connection options
type DatabaseConfig struct {
	Type         string `default:"mysql"`
	Host         string `default:"localhost"`
	Port         int    `default:"3306"`
	User         string `default:"floppy"`
	Password     string `default:"floppy"`
	DatabaseName string `default:"floppy"`
	MaxIdleConns int    `default:"5"`
	MaxOpenConns int    `default:"5"`
}

// OauthConfig stores oauth service configuration options
type OauthConfig struct {
	AccessTokenLifetime  int `default:"3600"`    // default to 1 hour
	RefreshTokenLifetime int `default:"1209600"` // default to 14 days
	AuthCodeLifetime     int `default:"3600"`    // default to 1 hour
}

// SessionConfig stores session configuration for the web app
type SessionConfig struct {
	Secret string `default:"test_secret"`
	Path   string `default:"/"`
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int `default:"604800"`
	// When you tag a cookie with the HttpOnly flag, it tells the browser that
	// this particular cookie should only be accessed by the server.
	// Any attempt to access the cookie from client script is strictly forbidden.
	HTTPOnly bool `default:"True"`
}

// Config stores all configuration options
type Config struct {
	Database DatabaseConfig
	Oauth    OauthConfig
	Session  SessionConfig
	//Seedfile	  []string	`default:nil`
	ServerPort    int  `default:"8080"`
	IsDevelopment bool `default:"True"`
	MinPwdLength  int  `default:"6"`
}

// DefaultConfig ...
// Let's start with some sensible defaults
var (
	ConfigInstance *Config
	defaultConfig  = &Config{
		Database: DatabaseConfig{
			Type:         "mysql",
			Host:         "localhost",
			Port:         3306,
			User:         "floppy_oauth2_server",
			Password:     "",
			DatabaseName: "floppy_oauth2_server",
			MaxIdleConns: 5,
			MaxOpenConns: 5,
		},
		Oauth: OauthConfig{
			AccessTokenLifetime:  3600,    // 1 hour
			RefreshTokenLifetime: 1209600, // 14 days
			AuthCodeLifetime:     3600,    // 1 hour
		},
		Session: SessionConfig{
			Secret:   "test_secret",
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
			HTTPOnly: true,
		},
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

		dbcfg := c.GetMap("/floppyisadog/foauthserver/database/")
		ConfigInstance.Database.Type = dbcfg["Type"]
		ConfigInstance.Database.Host = dbcfg["Host"]
		ConfigInstance.Database.Port = util.StringToInt(dbcfg["Port"])
		ConfigInstance.Database.User = dbcfg["User"]
		ConfigInstance.Database.Password = dbcfg["Password"]
		ConfigInstance.Database.DatabaseName = dbcfg["DatabaseName"]
		ConfigInstance.Database.MaxIdleConns = util.StringToInt(dbcfg["MaxIdleConns"])
		ConfigInstance.Database.MaxOpenConns = util.StringToInt(dbcfg["MaxOpenConns"])

		oauthcfg := c.GetMap("/floppyisadog/foauthserver/oauth/")
		ConfigInstance.Oauth.AccessTokenLifetime = util.StringToInt(oauthcfg["AccessTokenLifetime"])
		ConfigInstance.Oauth.RefreshTokenLifetime = util.StringToInt(oauthcfg["RefreshTokenLifetime"])
		ConfigInstance.Oauth.AuthCodeLifetime = util.StringToInt(oauthcfg["AuthCodeLifetime"])

		sessionCfg := c.GetMap("/floppyisadog/foauthserver/session/")
		ConfigInstance.Session.Secret = sessionCfg["Secret"]
		ConfigInstance.Session.Path = sessionCfg["Path"]
		ConfigInstance.Session.MaxAge = util.StringToInt(sessionCfg["MaxAge"])
		ConfigInstance.Session.HTTPOnly = util.StringToBool(sessionCfg["HTTPOnly"])

		ConfigInstance.ServerPort = c.GetIntWithDef("/floppyisadog/foauthserver/<ServerPort>", 8080)
		ConfigInstance.IsDevelopment = c.GetBoolWithDef("/floppyisadog/foauthserver/<IsDevelopment>", true)
		ConfigInstance.MinPwdLength = c.GetIntWithDef("/floppyisadog/foauthserver/<MinPwdLength>", 6)

	}
}

func GetConfig() *Config {
	if ConfigInstance == nil {
		return defaultConfig
	}

	return ConfigInstance
}
