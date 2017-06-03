package tonicdb

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lingmiaotech/tonic"
	"github.com/pressly/goose"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	command := "help"
	if len(args) > 0 {
		command = args[0]
	}

	var err error

	switch command {
	case "up":
		err = database("up", "./migrations")
	case "up-to":
		if len(args) < 2 {
			helpAndExit()
		}
		err = database("up-to", "./migrations", args[1])
	case "down":
		err = database("down", "./migrations")
	case "down-to":
		if len(args) < 2 {
			helpAndExit()
		}
		err = database("down-to", "./migrations", args[1])
	case "create":
		if len(args) < 2 {
			helpAndExit()
		}
		err = database("create", "./migrations", args[1], "sql")
	default:
		helpAndExit()
		return
	}

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
		return
	}

	log.Println("Done!")

}

func helpAndExit() {
	log.Print(
		`
		Avalable usages are:\n
		APP_ENV=./configs/development.yaml tonicdb help\n
		APP_ENV=./configs/development.yaml tonicdb up\n
		APP_ENV=./configs/development.yaml tonicdb up-to VERSION\n
		APP_ENV=./configs/development.yaml tonicdb down\n
		APP_ENV=./configs/development.yaml tonicdb down-to VERSION\n
		APP_ENV=./configs/development.yaml tonicdb create MIGRATION_NAME\n

		For production, run:\n
		APP_ENV=./configs/production.yaml tonicdb help\n
		APP_ENV=./configs/production.yaml tonicdb up\n
		APP_ENV=./configs/production.yaml tonicdb up-to VERSION\n
		APP_ENV=./configs/production.yaml tonicdb down\n
		APP_ENV=./configs/production.yaml tonicdb down-to VERSION\n
		APP_ENV=./configs/production.yaml tonicdb create MIGRATION_NAME\n
		`,
	)
	os.Exit(1)
}

func database(command string, dir string, args ...string) error {
	var err error

	err = tonic.InitConfigs()
	if err != nil {
		return err
	}

	driver := tonic.Configs.GetString("database.driver")
	if driver == "" {
		return errors.New("tonic_error.database.empty_driver_config")
	}

	err = goose.SetDialect(driver)
	if err != nil {
		return err
	}

	dbstring := tonic.Configs.GetString("database.dbstring")
	if dbstring == "" {
		return errors.New("tonic_error.database.empty_dbstring_config")
	}

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		return err
	}

	err = goose.Run(command, db, dir, args...)
	if err != nil {
		return err
	}

	return nil
}
