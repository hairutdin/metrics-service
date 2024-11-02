package storage

import (
	"context"
	"fmt"

	"github.com/hairutdin/metrics-service/models"

	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	DB *pgx.Conn
}

func NewPostgresStorage(dbConn *pgx.Conn) *PostgresStorage {
	return &PostgresStorage{DB: dbConn}
}

func (s *PostgresStorage) UpdateGauge(name string, value float64) {
	query := `
		INSERT INTO gauge_metrics (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
	`

	_, err := s.DB.Exec(context.Background(), query, name, value)
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
	_, err := s.DB.Exec(context.Background(), query, name, delta)
	if err != nil {
		fmt.Printf("Error updating counter metric: %v\n", err)
	}
}

func (s *PostgresStorage) UpdateMetricsBatch(metrics []models.Metrics) error {
	tx, err := s.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, metric := range metrics {
		if metric.MType == "gauge" && metric.Value != nil {
			_, err := tx.Exec(context.Background(),
				`INSERT INTO gauge_metrics (name, value) VALUES ($1, $2)
				 ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value`, metric.ID, *metric.Value)
			if err != nil {
				return err
			}
		} else if metric.MType == "counter" && metric.Delta != nil {
			_, err := tx.Exec(context.Background(),
				`INSERT INTO counter_metrics (name, value) VALUES ($1, $2)
				 ON CONFLICT (name) DO UPDATE SET value = counter_metrics.value + EXCLUDED.value`, metric.ID, *metric.Delta)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(context.Background())
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

	err := s.DB.QueryRow(context.Background(), query, name).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("metric not found: %w", err)
	}

	return result, nil
}

func (s *PostgresStorage) GetAllMetrics() map[string]string {
	metrics := make(map[string]string)

	gaugeQuery := `SELECT name, value FROM gauge_metrics`
	counterQuery := `SELECT name, value FROM counter_metrics`

	rows, err := s.DB.Query(context.Background(), gaugeQuery)
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

	rows, err = s.DB.Query(context.Background(), counterQuery)
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
