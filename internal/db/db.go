package db

import (
	"context"

	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

func ConnectToDB(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func InitializeTables(conn *pgx.Conn) error {
	gaugeTableQuery := `
	CREATE TABLE IF NOT EXISTS gauge_metrics (
		name TEXT PRIMARY KEY,
		value DOUBLE PRECISION
	);
	`
	counterTableQuery := `
	CREATE TABLE IF NOT EXISTS counter_metrics (
		name TEXT PRIMARY KEY,
		value BIGINT
	)
	`

	_, err := conn.Exec(context.Background(), gaugeTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create gauge_metrics table: %w", err)
	}

	_, err = conn.Exec(context.Background(), counterTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create counter_metrics table: %w", err)
	}

	return nil
}

func PingDB(conn *pgx.Conn) error {
	return conn.Ping(context.Background())
}

func CloseDB(conn *pgx.Conn) {
	if conn != nil {
		conn.Close(context.Background())
	}
}
