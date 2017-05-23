package templates

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lingmiaotech/aviation"
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
		APP_ENV=./configs/development.yaml avdb help\n
		APP_ENV=./configs/development.yaml avdb up\n
		APP_ENV=./configs/development.yaml avdb up-to VERSION\n
		APP_ENV=./configs/development.yaml avdb down\n
		APP_ENV=./configs/development.yaml avdb down-to VERSION\n
		APP_ENV=./configs/development.yaml avdb create MIGRATION_NAME\n

		For production, run:\n
		APP_ENV=./configs/production.yaml avdb help\n
		APP_ENV=./configs/production.yaml avdb up\n
		APP_ENV=./configs/production.yaml avdb up-to VERSION\n
		APP_ENV=./configs/production.yaml avdb down\n
		APP_ENV=./configs/production.yaml avdb down-to VERSION\n
		APP_ENV=./configs/production.yaml avdb create MIGRATION_NAME\n
		`,
	)
	os.Exit(1)
}

func database(command string, dir string, args ...string) error {
	var err error

	err = aviation.InitConfigs()
	if err != nil {
		return err
	}

	driver := aviation.Configs.GetString("database.driver")
	if driver == "" {
		return errors.New("aviation_error.database.empty_driver_config")
	}

	err = goose.SetDialect(driver)
	if err != nil {
		return err
	}

	dbstring := aviation.Configs.GetString("database.dbstring")
	if dbstring == "" {
		return errors.New("aviation_error.database.empty_dbstring_config")
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
