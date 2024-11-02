package storage

import (
	"context"
	"testing"

	"github.com/hairutdin/metrics-service/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type MockBatchResults struct{}

func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *MockBatchResults) Query() (pgx.Rows, error) {
	return nil, nil
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	return nil
}

func (m *MockBatchResults) Close() error {
	return nil
}

type MockTx struct {
	ExecCount   *int
	ShouldRetry bool
}

func (m *MockTx) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	*m.ExecCount++
	if m.ShouldRetry && *m.ExecCount <= 3 {
		return pgconn.CommandTag{}, &pgconn.PgError{Code: pgerrcode.ConnectionException}
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	return nil
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (m *MockTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *MockTx) Commit(ctx context.Context) error {
	if m.ShouldRetry && *m.ExecCount <= 3 {
		return &pgconn.PgError{Code: pgerrcode.ConnectionException}
	}
	return nil
}

func (m *MockTx) Conn() *pgx.Conn {
	return nil
}

func (m *MockTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *MockTx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *MockTx) QueryRow(context.Context, string, ...any) pgx.Row {
	return nil
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *MockTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults {
	return &MockBatchResults{}
}

type MockPostgresStorageConn struct {
	ExecCount   int
	ShouldRetry bool
}

func (m *MockPostgresStorageConn) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	m.ExecCount++
	if m.ShouldRetry && m.ExecCount <= 3 {
		return pgconn.CommandTag{}, &pgconn.PgError{Code: pgerrcode.ConnectionException}
	}
	return pgconn.CommandTag{}, nil
}

func (m *MockPostgresStorageConn) Begin(ctx context.Context) (pgx.Tx, error) {
	return &MockTx{ExecCount: &m.ExecCount, ShouldRetry: m.ShouldRetry}, nil
}

func (m *MockPostgresStorageConn) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return nil
}

func (m *MockPostgresStorageConn) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func TestUpdateMetricsBatch_RetryOnDatabaseError(t *testing.T) {
	mockConn := &MockPostgresStorageConn{ShouldRetry: true}
	storage := NewPostgresStorage(mockConn)

	metrics := []models.Metrics{
		{ID: "test_gauge", MType: "gauge", Value: func(v float64) *float64 { return &v }(42.0)},
	}

	err := storage.UpdateMetricsBatch(metrics)

	assert.NoError(t, err, "Expected no error due to retry mechanism")
	assert.Equal(t, 4, mockConn.ExecCount, "Expected 4 attempts: 1 initial + 3 retries")
}
