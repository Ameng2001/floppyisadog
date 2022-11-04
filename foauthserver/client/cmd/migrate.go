package cmd

import (
	"github.com/floppyisadog/appcommon/utils/database"
	"github.com/floppyisadog/foauthserver/managers/configmgr"
	"github.com/floppyisadog/foauthserver/models/migrations"
)

// Migrate runs database migrations
func Migrate(configFile string) error {
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

	// Bootstrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := migrations.MigrateAllTables(db); err != nil {
		return err
	}

	return nil
}
