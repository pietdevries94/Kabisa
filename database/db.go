package database

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog"

	// the sql driver
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

// init creates db connection and runs migrations
func Init(logger *zerolog.Logger, dsn string) *sql.DB {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		logger.Fatal().Err(err).Str("DSN", dsn).Msg("could not open sqlite db")
	}
	// This is to prevent the db from having locks. This is a sqlite specific issue
	db.SetMaxOpenConns(1)

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Fatal().Err(err).Str("DSN", dsn).Msg("could not create source for migration")
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		logger.Fatal().Err(err).Str("DSN", dsn).Msg("could not create instance for migration")
	}

	m, err := migrate.NewWithInstance(
		"migrations", source,
		"sqlite", driver,
	)
	if err != nil {
		logger.Fatal().Err(err).Str("DSN", dsn).Msg("could not create migrator with source")
	}
	err = m.Up()
	if err != nil {
		logger.Fatal().Err(err).Str("DSN", dsn).Msg("could not run migration")
	}

	return db
}
