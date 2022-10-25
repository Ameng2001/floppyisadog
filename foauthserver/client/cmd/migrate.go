package cmd

import (
	"github.com/floppyisadog/foauthserver/util/config"
	"github.com/floppyisadog/foauthserver/util/database"
	"github.com/floppyisadog/foauthserver/util/migrations"
)

// Migrate runs database migrations
func Migrate(configFile string) error {
	config.InitConfig(configFile)

	db, err := database.InitDB(config.GetConfig())
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
