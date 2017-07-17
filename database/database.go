package database

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lingmiaotech/tonic/configs"
	"time"
)

var Engine *gorm.DB

func InitDatabase() (err error) {
	isDatabaseSet := configs.IsSet("database")
	if !isDatabaseSet {
		return nil
	}

	driver := configs.GetString("database.driver")
	if driver == "" {
		return errors.New("tonic_error.database.empty_dbstring_config")
	}

	dbstring := configs.GetString("database.dbstring")
	if dbstring == "" {
		return errors.New("tonic_error.database.empty_dbstring_config")
	}

	Engine, err = gorm.Open(driver, dbstring)
	if err != nil {
		return
	}

	Engine.DB().SetConnMaxLifetime(1 * time.Hour)
	Engine.DB().SetMaxOpenConns(15)
	return
}
