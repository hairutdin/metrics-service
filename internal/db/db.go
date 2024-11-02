package db

import (
	"context"

	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

var DB *pgx.Conn

func ConnectToDB(dsn string) error {
	var err error
	DB, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}

	err = initializeTables()
	if err != nil {
		return err
	}

	return nil
}

func initializeTables() error {
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

	_, err := DB.Exec(context.Background(), gaugeTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create gauge_metrics table: %w", err)
	}

	_, err = DB.Exec(context.Background(), counterTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create counter_metrics table: %w", err)
	}

	return nil
}

func PingDB() error {
	return DB.Ping(context.Background())
}

func CloseDB() {
	if DB != nil {
		DB.Close(context.Background())
	}
}
