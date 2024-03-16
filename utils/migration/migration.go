package migration

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Conf struct {
	DB             *sql.DB
	Path           string
	MigrationTable string
	DatabaseName   string
}

func New(m Conf) (*migrate.Migrate, error) {
	pConf := postgres.Config{}
	if m.MigrationTable != "" {
		pConf.MigrationsTable = m.MigrationTable
	}
	if m.DatabaseName != "" {
		pConf.DatabaseName = m.DatabaseName
	}
	driver, err := postgres.WithInstance(m.DB, &pConf)
	if err != nil {
		return &migrate.Migrate{}, err
	}
	instance, err := migrate.NewWithDatabaseInstance(
		m.Path,
		"postgres", driver)
	if err != nil {
		return &migrate.Migrate{}, err
	}
	return instance, err
}
