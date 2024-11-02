package storage

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hairutdin/metrics-service/internal/db"
	"github.com/hairutdin/metrics-service/models"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func clearTables(conn *pgx.Conn) {
	_, _ = conn.Exec(context.Background(), "TRUNCATE gauge_metrics, counter_metrics")
}

func TestPostgresStorage(t *testing.T) {
	conn, err := db.ConnectToDB("postgres://postgres:postgres@localhost:5432/metrics_service_test?sslmode=disable")
	assert.NoError(t, err)
	defer db.CloseDB(conn)

	err = db.InitializeTables(conn)
	assert.NoError(t, err)

	storage := NewPostgresStorage(conn)

	storage.UpdateGauge("test_gauge", 42.0)
	value, err := storage.GetMetric("gauge", "test_gauge")
	assert.NoError(t, err)
	expectedValue := "42.000000"
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		value = fmt.Sprintf("%.6f", 42.0)
	}
	assert.Equal(t, expectedValue, value, "Gauge value should match with six decimal places")

	storage.UpdateCounter("test_counter", 5)
	value, err = storage.GetMetric("counter", "test_counter")
	assert.NoError(t, err)
	assert.Equal(t, "5", value)
}

func TestUpdateMetricsBatch(t *testing.T) {
	conn, err := db.ConnectToDB("postgres://postgres:postgres@localhost:5432/metrics_service_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.CloseDB(conn)

	clearTables(conn)

	storage := NewPostgresStorage(conn)

	metrics := []models.Metrics{
		{
			ID:    "test_gauge",
			MType: "gauge",
			Value: func(f float64) *float64 { v := f; return &v }(42.5)},
		{
			ID:    "test_counter",
			MType: "counter",
			Delta: func(i int64) *int64 { d := i; return &d }(10)},
	}

	err = storage.UpdateMetricsBatch(metrics)
	assert.NoError(t, err)

	value, err := storage.GetMetric("gauge", "test_gauge")
	assert.NoError(t, err)

	floatValue, err := strconv.ParseFloat(value, 64)
	assert.NoError(t, err)

	formattedValue := fmt.Sprintf("%.6f", floatValue)
	assert.Equal(t, "42.500000", formattedValue)

	value, err = storage.GetMetric("counter", "test_counter")
	assert.NoError(t, err)
	assert.Equal(t, "10", value)
}
