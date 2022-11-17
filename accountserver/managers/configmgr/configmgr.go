package configmgr

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars/util/conf"
	"github.com/floppyisadog/appcommon/utils"
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

// Config stores all configuration options
type Config struct {
	Database      DatabaseConfig
	Outerfactory  map[string]string
	IsDevelopment bool `default:"True"`
	SigningToken  string
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
		IsDevelopment: true,
		SigningToken:  "joystaff-account",
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

		dbcfg := c.GetMap("/floppyisadog/accountserver/database/")
		ConfigInstance.Database.Type = dbcfg["Type"]
		ConfigInstance.Database.Host = dbcfg["Host"]
		ConfigInstance.Database.Port = utils.StringToInt(dbcfg["Port"])
		ConfigInstance.Database.User = dbcfg["User"]
		ConfigInstance.Database.Password = dbcfg["Password"]
		ConfigInstance.Database.DatabaseName = dbcfg["DatabaseName"]
		ConfigInstance.Database.MaxIdleConns = utils.StringToInt(dbcfg["MaxIdleConns"])
		ConfigInstance.Database.MaxOpenConns = utils.StringToInt(dbcfg["MaxOpenConns"])

		ConfigInstance.Outerfactory = c.GetMap("/floppyisadog/accountserver/outerfactory/")

		ConfigInstance.IsDevelopment = c.GetBoolWithDef("/floppyisadog/accountserver/<IsDevelopment>", true)
		ConfigInstance.SigningToken = c.GetString("/floppyisadog/accountserver/<SigningToken>")
	}
}

func GetConfig() *Config {
	if ConfigInstance == nil {
		return defaultConfig
	}

	return ConfigInstance
}
