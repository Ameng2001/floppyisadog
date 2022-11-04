package cmd

import (
	// "errors"
	// "strings"

	"github.com/RichardKnop/go-fixtures"
	"github.com/floppyisadog/appcommon/utils/database"
	"github.com/floppyisadog/foauthserver/managers/configmgr"
)

// LoadData loads fixtures
// func LoadDataHelper(command string) (string,error) {
// 	var msg string
// 	cmd := strings.Split(command, " ")
// 	if len(cmd) <= 1 {
// 		msg = "loaddata param error see help"
// 		return errors.New("loaddata param error see help")
// 	}

// 	if cmd[0] != "loaddata" {
// 		msg = "unexpected command"
// 		return msg, errors.New(msg)
// 	}

// 	configFile := "foauthserver.conf"
// 	var paths []string

// 	cfg, db, err := initConfigDB(configFile)
// 	if err != nil {
// 		return err
// 	}

// 	for _, v := range cfg.Seedfile {
// 		paths = append(paths,v)
// 	}

// 	defer db.Close()
// 	return fixtures.LoadFiles(paths, db.DB(), cfg.Database.Type)
// }

func LoadData(paths []string, configFile string) error {
	configmgr.InitConfig(configFile)

	dbcf := configmgr.GetConfig().Database
	db, err := database.InitDB(dbcf.Type,
		dbcf.User,
		dbcf.Password,
		dbcf.Host,
		dbcf.DatabaseName,
		dbcf.Port,
		dbcf.MaxIdleConns,
		dbcf.MaxOpenConns,
		configmgr.GetConfig().IsDevelopment)
	if err != nil {
		return err
	}

	defer db.Close()
	return fixtures.LoadFiles(paths, db.DB(), dbcf.Type)
}
