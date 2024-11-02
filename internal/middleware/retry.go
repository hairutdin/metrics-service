package middleware

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var retryIntervals = []time.Duration{
	time.Second,
	3 * time.Second,
	5 * time.Second,
}

var maxRetries = len(retryIntervals)

func RetryOperation(operation func() error) error {
	for i, interval := range retryIntervals {
		err := operation()
		if err == nil {
			return nil
		}
		if !isRetriable(err) {
			return err
		}
		log.Printf("Attempt %d failed: %v. Retrying in %v...", i+1, err, interval)
		time.Sleep(interval)
	}

	err := operation()
	if err == nil {
		return nil
	}

	return fmt.Errorf("operation failed after %d retries", maxRetries)
}

func isRetriable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.ConnectionException
	}

	return false
}
