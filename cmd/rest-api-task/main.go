package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/ihrk/rest-api-task/internal/app"
	"github.com/ihrk/rest-api-task/internal/app/database/migrations"
)

const defaultConfigPath = "./configs/config.yml"

func main() {
	cfgPathPtr := flag.String("config", defaultConfigPath, "path to app config file")
	migratePtr := flag.String("migrate", "", "run db migration, allowed values: 'up', 'down'")

	flag.Parse()

	if direction := *migratePtr; direction != "" {
		migrate(*cfgPathPtr, direction)

		return
	}

	log.Fatal(app.Run(*cfgPathPtr))
}

func migrate(cfgPath, direction string) {
	cfg := app.LoadConfig(cfgPath)

	m, err := migrations.NewMigrate(cfg.DB.URL())
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		err = fmt.Errorf("unknown migrate direction: %s", direction)
	}

	if err != nil {
		log.Fatal(err)
	}
}
