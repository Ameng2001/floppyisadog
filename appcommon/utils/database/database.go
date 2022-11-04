package database

import (
	"fmt"
	"time"

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

func InitDB(Type, User, Password, Host, DatabaseName string, Port, MaxIdleConns, MaxOpenConns int, IsDevelopment bool) (*gorm.DB, error) {
	// Mysql
	if Type == "mysql" {
		// Connection args
		// see https://github.com/go-sql-driver/mysql#dsn-data-source-name
		args := fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			User,
			Password,
			Host,
			Port,
			DatabaseName,
		)

		var err error
		DB, err = gorm.Open(Type, args)
		if err != nil {
			return DB, err
		}

		// Max idle connections
		DB.DB().SetMaxIdleConns(MaxIdleConns)

		// Max open connections
		DB.DB().SetMaxOpenConns(MaxOpenConns)

		// Database logging
		DB.LogMode(IsDevelopment)

		return DB, nil
	}

	// Database type not supported
	return nil, fmt.Errorf("database type %s not supported", Type)
}
