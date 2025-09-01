package db

import (
	"crud/internal/config"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func Connect(cfg *config.Config) *bun.DB {
	connector := pgdriver.NewConnector(pgdriver.WithDSN(cfg.PostgresDSN))
	sqldb := sql.OpenDB(connector)

	// Production-grade pool settings
	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	return bun.NewDB(sqldb, pgdialect.New())
}
