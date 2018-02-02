package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lingmiaotech/tonic/configs"
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

	appName := configs.GetString("app_name")
	username := configs.GetString("database.username")
	password := configs.GetDynamicString("database.password")
	host := configs.GetString("database.host")
	port := configs.GetString("database.port")
	args := configs.GetString("database.args")

	dbstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", username, password, host, port, appName, args)
	Engine, err = gorm.Open(driver, dbstring)
	if err != nil {
		return
	}

	Engine.DB().SetConnMaxLifetime(1 * time.Hour)
	Engine.DB().SetMaxOpenConns(15)
	return
}
