package migrations

import (
	"fmt"

	"github.com/floppyisadog/foauthserver/models"
	"github.com/jinzhu/gorm"
)

// Define the migration stage
type MigrationStage struct {
	Name     string
	Function func(db *gorm.DB) error
}

// A single database migration
type Migration struct {
	gorm.Model
	Name string `sql:"size:255"`
}

// init tables
var (
	list = []MigrationStage{
		{
			Name:     "initial",
			Function: migrate0001,
		},
	}
)

// ///////////////////////////////////////////////////
func Bootstrap(db *gorm.DB) error {
	migrationName := "bootstrap_migrations"

	mig := new(Migration)
	mig.Name = migrationName

	exists := nil == db.Where("name = ?", migrationName).First(mig).Error
	if exists {
		fmt.Printf("Skipping %s migration", migrationName)
		return nil
	}

	fmt.Printf("Begin %s migration", migrationName)

	if err := db.CreateTable(new(Migration)).Error; err != nil {
		return fmt.Errorf("error creating migration table %s", db.Error)
	}

	if err := db.Create(mig).Error; err != nil {
		return fmt.Errorf("error saving record to migration table:%s", err)
	}

	return nil
}

func MigrateAllTables(db *gorm.DB) error {
	return migrate(db, list)
}

// migrate functions
func migrate(db *gorm.DB, migrations []MigrationStage) error {
	for _, m := range migrations {
		if migrationExists(db, m.Name) {
			continue
		}

		if err := m.Function(db); err != nil {
			return err
		}

		if err := saveMigration(db, m.Name); err != nil {
			return err
		}
	}

	return nil
}

func migrationExists(db *gorm.DB, migrationName string) bool {
	mig := new(Migration)
	found := !db.Where("name = ?", migrationName).First(mig).RecordNotFound()

	if found {
		fmt.Printf("Skipping %s migration", migrationName)
	} else {
		fmt.Printf("Running %s migration", migrationName)
	}

	return found
}

func saveMigration(db *gorm.DB, migrationName string) error {
	mig := new(Migration)
	mig.Name = migrationName

	if err := db.Create(mig).Error; err != nil {
		fmt.Printf("Error save record to migration table %s", err)
		return fmt.Errorf("error save record to migration table %s", err)
	}

	return nil
}

func migrate0001(db *gorm.DB) error {
	// Create tables
	if err := db.CreateTable(new(models.OauthClient)).Error; err != nil {
		return fmt.Errorf("error creating oauth_clients table: %s", err)
	}
	if err := db.CreateTable(new(models.OauthScope)).Error; err != nil {
		return fmt.Errorf("error creating oauth_scopes table: %s", err)
	}
	if err := db.CreateTable(new(models.OauthRole)).Error; err != nil {
		return fmt.Errorf("error creating oauth_roles table: %s", err)
	}
	if err := db.CreateTable(new(models.OauthUser)).Error; err != nil {
		return fmt.Errorf("error creating oauth_users table: %s", err)
	}
	if err := db.CreateTable(new(models.OauthRefreshToken)).Error; err != nil {
		return fmt.Errorf("error creating oauth_refresh_tokens table: %s", err)
	}
	if err := db.CreateTable(new(models.OauthAccessToken)).Error; err != nil {
		return fmt.Errorf("error creating oauth_access_tokens table: %s", err)
	}
	if err := db.CreateTable(new(models.OauthAuthorizationCode)).Error; err != nil {
		return fmt.Errorf("error creating oauth_authorization_codes table: %s", err)
	}
	err := db.Model(new(models.OauthUser)).AddForeignKey(
		"role_id", "oauth_roles(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_users.role_id for oauth_roles(id): %s", err)
	}
	err = db.Model(new(models.OauthRefreshToken)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_refresh_tokens.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(models.OauthRefreshToken)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_refresh_tokens.user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(models.OauthAccessToken)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_access_tokens.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(models.OauthAccessToken)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_access_tokens.user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(models.OauthAuthorizationCode)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_authorization_codes.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(models.OauthAuthorizationCode)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("error creating foreign key on "+
			"oauth_authorization_codes.user_id for oauth_users(id): %s", err)
	}

	return nil
}
