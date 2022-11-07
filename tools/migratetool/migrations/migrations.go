package migrations

import (
	"fmt"

	accountmodels "github.com/floppyisadog/accountserver/models"
	companymodels "github.com/floppyisadog/companyserver/models"
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
	if err := db.CreateTable(new(accountmodels.Account)).Error; err != nil {
		return fmt.Errorf("error creating accounts table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Company)).Error; err != nil {
		return fmt.Errorf("error creating companies table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Admin)).Error; err != nil {
		return fmt.Errorf("error creating admins table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Team)).Error; err != nil {
		return fmt.Errorf("error creating teams table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Directory)).Error; err != nil {
		return fmt.Errorf("error creating directories table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Worker)).Error; err != nil {
		return fmt.Errorf("error creating workers table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Job)).Error; err != nil {
		return fmt.Errorf("error creating jobs table: %s", err)
	}
	if err := db.CreateTable(new(companymodels.Shift)).Error; err != nil {
		return fmt.Errorf("error creating shifts table: %s", err)
	}

	return nil
}
