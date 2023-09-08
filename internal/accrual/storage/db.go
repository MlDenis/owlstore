package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

// инизиацлизация бд
func InitDB(configDB string, migratePath string, log *zap.Logger) (*PostgresDB, error) {
	// err := runMigrations(configDB, migratePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	// }
	pool, err := pgxpool.New(context.Background(), configDB)
	if err == nil {
		return &PostgresDB{pool: pool}, nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		log.Error("Database initialization error", zap.Error(err))
		pool, err := pgxpool.New(context.Background(), configDB)
		if err == nil {
			log.Info("Successful database connection")
			return &PostgresDB{pool: pool}, nil
		}

	}
	return nil, fmt.Errorf("failed to create a connection pool: %w", err)
}

// функция чтобы закрыть соедининение
func (pgdb *PostgresDB) Close() {
	pgdb.pool.Close()
}

// накатываем миграции
func runMigrations(dsn string, migratePath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migratePath), dsn)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
