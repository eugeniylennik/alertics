package database

import (
	"context"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
}

func NewClient(ctx context.Context, maxAttempt int, dsn string) (conn *pgxpool.Pool, err error) {
	if err := utils.RetryConnectToDataBase(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		conn, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempt, 5*time.Second); err != nil {
		log.Fatalf("failed to retry connection to database, %s", err)
	}
	if err := createTable(conn); err != nil {
		log.Fatalf("failed to create table metrics, %s", err)
	}
	return
}

func createTable(conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS metrics (
            id TEXT PRIMARY KEY,
            type TEXT NOT NULL,
            delta BIGINT,
            value DOUBLE PRECISION,
            hash TEXT
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
