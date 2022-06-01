package migrations

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var embeddedFS embed.FS

func NewMigrate(databaseURL string) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(embeddedFS, ".")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithSourceInstance("file", sourceDriver, databaseURL)
	if err != nil {
		return nil, err
	}

	return m, nil
}
