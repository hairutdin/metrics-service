package storage

import (
	"context"
	"fmt"

	"github.com/hairutdin/metrics-service/internal/db"
)

type PostgresStorage struct{}

func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{}
}

func (s *PostgresStorage) UpdateGauge(name string, value float64) {
	query := `
		INSERT INTO gauge_metrics (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
	`

	_, err := db.DB.Exec(context.Background(), query, name, value)
	if err != nil {
		fmt.Printf("Error updating gauge metric: %v\n", err)
	}
}

func (s *PostgresStorage) UpdateCounter(name string, delta int64) {
	query := `
		INSERT INTO counter_metrics (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;
	`
	_, err := db.DB.Exec(context.Background(), query, name, delta)
	if err != nil {
		fmt.Printf("Error updating counter metric: %v\n", err)
	}
}

func (s *PostgresStorage) GetMetric(metricType, name string) (string, error) {
	var query string
	var result string

	switch metricType {
	case "gauge":
		query = "SELECT value FROM gauge_metrics WHERE name = $1"
	case "counter":
		query = "SELECT value FROM counter_metrics WHERE name = $1"
	default:
		return "", fmt.Errorf("invalid metric type: %s", metricType)
	}

	err := db.DB.QueryRow(context.Background(), query, name).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("metric not found: %w", err)
	}

	return result, nil
}

func (s *PostgresStorage) GetAllMetrics() map[string]string {
	metrics := make(map[string]string)

	gaugeQuery := `SELECT name, value FROM gauge_metrics`
	counterQuery := `SELECT name, value FROM counter_metrics`

	rows, err := db.DB.Query(context.Background(), gaugeQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var value float64
			if err := rows.Scan(&name, &value); err == nil {
				metrics[name] = fmt.Sprintf("gauge: %f", value)
			}
		}
	}

	rows, err = db.DB.Query(context.Background(), counterQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var value int64
			if err := rows.Scan(&name, &value); err == nil {
				metrics[name] = fmt.Sprintf("counter: %d", value)
			}
		}
	}

	return metrics

}
