package tonic

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var Database *gorm.DB

func InitDatabase() (err error) {
	isDatabaseSet := Configs.IsSet("database")
	if !isDatabaseSet {
		return nil
	}

	driver := Configs.GetString("database.driver")
	if driver == "" {
		return errors.New("tonic_error.database.empty_dbstring_config")
	}

	dbstring := Configs.GetString("database.dbstring")
	if dbstring == "" {
		return errors.New("tonic_error.database.empty_dbstring_config")
	}

	Database, err = gorm.Open(driver, dbstring)
	if err != nil {
		return
	}

	Database.DB().SetConnMaxLifetime(1 * time.Hour)
	Database.DB().SetMaxOpenConns(15)
	return
}
