package database

import (
	"fmt"
	"time"

	"github.com/floppyisadog/foauthserver/util/config"
	"github.com/jinzhu/gorm"

	// Driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}

// func init() {
// 	gorm.Nowfunc = func() time.Time {
// 		return time.Now().UTC()
// 	}
// }

func InitDB(cf *config.Config) (*gorm.DB, error) {
	// Mysql
	if cf.Database.Type == "mysql" {
		// Connection args
		// see https://github.com/go-sql-driver/mysql#dsn-data-source-name
		args := fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			cf.Database.User,
			cf.Database.Password,
			cf.Database.Host,
			cf.Database.Port,
			cf.Database.DatabaseName,
		)

		var err error
		DB, err = gorm.Open(cf.Database.Type, args)
		if err != nil {
			return DB, err
		}

		// Max idle connections
		DB.DB().SetMaxIdleConns(cf.Database.MaxIdleConns)

		// Max open connections
		DB.DB().SetMaxOpenConns(cf.Database.MaxOpenConns)

		// Database logging
		DB.LogMode(cf.IsDevelopment)

		return DB, nil
	}

	// Database type not supported
	return nil, fmt.Errorf("database type %s not supported", cf.Database.Type)
}
