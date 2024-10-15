package db

import (
	"context"

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
	return nil
}

func PingDB() error {
	return DB.Ping(context.Background())
}

func CloseDB() {
	DB.Close(context.Background())
}
