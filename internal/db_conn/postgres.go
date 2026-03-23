package db_conn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Fankemp/GameMatch/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Драйвер для работы с Postgres в миграциях
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Драйвер для чтения файлов миграций
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DB  *sqlx.DB
	dsn string
}

func NewDB(cfg *config.PostgreConfig) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlxDB, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, err
	}

	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(25)
	sqlxDB.SetConnMaxLifetime(5 * time.Minute)

	err = sqlxDB.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		DB:  sqlxDB,
		dsn: dsn,
	}, nil
}

func (d *DB) Migrate() error {
	sourceURL := "file://internal/db_conn/migrations"

	m, err := migrate.New(sourceURL, d.dsn)
	if err != nil {
		return fmt.Errorf("failed to migrate:  %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no new migration")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (d *DB) MigrateDown() error {
	sourceURL := "file://internal/db_conn/migration"
	m, err := migrate.New(sourceURL, d.dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err = m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("no migration to rollback")
		}

		return fmt.Errorf("failed to rollback migation: %w", err)
	}

	log.Println(" Migration rolled back")
	return nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}
