package db

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestNewPGXTracker(t *testing.T) {
	ctx := context.Background()
	host := getEnvOrDefault("PG_HOST", "localhost")
	port := getEnvInt("PG_PORT", 5432)
	database := getEnvOrDefault("PG_DATABASE", "testdb")

	tracker, err := NewPGXTracker(ctx, host, port, database)
	require.NoError(t, err)
	require.NotNil(t, tracker)

	assert.NotNil(t, tracker.Metrics)
	assert.NotNil(t, tracker.tracer)
	assert.NotEmpty(t, tracker.commonTraceAttrs)
}

func TestPGXTracker_TraceQueryStart(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	data := pgx.TraceQueryStartData{
		SQL: "SELECT * FROM users WHERE id = $1",
	}

	newCtx := tracker.TraceQueryStart(ctx, nil, data)
	assert.NotEqual(t, ctx, newCtx)

	// Verify metadata is stored in context
	meta, ok := newCtx.Value(pgQueryMetadataKey{}).(pgQueryMetadata)
	assert.True(t, ok)
	assert.Equal(t, "SELECT users", meta.spanName)
	assert.Contains(t, meta.attrs, semconv.DBQueryText(data.SQL))
}

func TestPGXTracker_TraceQueryEnd(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// First start a query to set up context
	startData := pgx.TraceQueryStartData{SQL: "SELECT * FROM users"}
	ctxWithMeta := tracker.TraceQueryStart(ctx, nil, startData)

	// Then end it
	commandTag := pgconn.NewCommandTag("SELECT 2")
	endData := pgx.TraceQueryEndData{
		CommandTag: commandTag,
		Err:        nil,
	}

	tracker.TraceQueryEnd(ctxWithMeta, nil, endData)

	// Should not panic
}

func TestPGXTracker_TraceQueryEndWithError(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// First start a query
	startData := pgx.TraceQueryStartData{SQL: "SELECT * FROM users"}
	ctxWithMeta := tracker.TraceQueryStart(ctx, nil, startData)

	// End with error
	pgErr := &pgconn.PgError{
		Code: pgerrcode.ConnectionException,
	}
	endData := pgx.TraceQueryEndData{
		CommandTag: pgconn.CommandTag{},
		Err:        pgErr,
	}

	tracker.TraceQueryEnd(ctxWithMeta, nil, endData)

	// Should not panic
}

func TestPGXTracker_TraceBatchStart(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	batch := &pgx.Batch{}
	batch.Queue("SELECT 1")
	batch.Queue("SELECT 2")

	data := pgx.TraceBatchStartData{Batch: batch}

	newCtx := tracker.TraceBatchStart(ctx, nil, data)
	assert.NotEqual(t, ctx, newCtx)

	// Verify batch size is stored
	batchSize, ok := newCtx.Value(batchSizeKey{}).(int)
	assert.True(t, ok)
	assert.Equal(t, 2, batchSize)
}

func TestPGXTracker_TraceBatchQuery(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	data := pgx.TraceBatchQueryData{
		SQL: "SELECT * FROM products",
	}

	tracker.TraceBatchQuery(ctx, nil, data)

	// Should not panic
}

func TestPGXTracker_TraceBatchQueryWithSpan(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Start a span context to test the span recording path
	ctx, span := tracker.tracer.Start(ctx, "test-span")
	defer span.End()

	data := pgx.TraceBatchQueryData{
		SQL: "SELECT * FROM products",
	}

	tracker.TraceBatchQuery(ctx, nil, data)

	// Should not panic and should add event to span
}

func TestPGXTracker_TraceBatchEnd(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Start a batch first
	batch := &pgx.Batch{}
	batch.Queue("SELECT 1")
	startData := pgx.TraceBatchStartData{Batch: batch}
	ctxWithMeta := tracker.TraceBatchStart(ctx, nil, startData)

	// End the batch
	endData := pgx.TraceBatchEndData{Err: nil}
	tracker.TraceBatchEnd(ctxWithMeta, nil, endData)

	// Should not panic
}

func TestPGXTracker_TraceBatchEndWithError(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Start a batch first
	batch := &pgx.Batch{}
	batch.Queue("SELECT 1")
	startData := pgx.TraceBatchStartData{Batch: batch}
	ctxWithMeta := tracker.TraceBatchStart(ctx, nil, startData)

	// End the batch with error
	pgErr := &pgconn.PgError{Code: pgerrcode.ConnectionException}
	endData := pgx.TraceBatchEndData{Err: pgErr}
	tracker.TraceBatchEnd(ctxWithMeta, nil, endData)

	// Should not panic
}

func TestPGXTracker_TraceBatchEndWithBatchSizeInSpan(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Create a batch with multiple queries
	batch := &pgx.Batch{}
	batch.Queue("SELECT 1")
	batch.Queue("SELECT 2")
	batch.Queue("SELECT 3")
	startData := pgx.TraceBatchStartData{Batch: batch}

	// Start the batch - this stores batch size in context
	ctxWithMeta := tracker.TraceBatchStart(ctx, nil, startData)

	// Verify batch size is in context
	batchSize, ok := ctxWithMeta.Value(batchSizeKey{}).(int)
	assert.True(t, ok)
	assert.Equal(t, 3, batchSize)

	// End the batch successfully
	endData := pgx.TraceBatchEndData{Err: nil}
	tracker.TraceBatchEnd(ctxWithMeta, nil, endData)

	// The batch size attribute should have been added to the span
	// (We can't directly verify the span contents in this test, but we can
	// verify the batch size was retrieved from context correctly)
}

func TestPGXTracker_TraceConnectStart(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	data := pgx.TraceConnectStartData{}

	newCtx := tracker.TraceConnectStart(ctx, data)
	assert.NotEqual(t, ctx, newCtx)
}

func TestPGXTracker_TraceConnectEnd(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Start connect first
	ctxWithMeta := tracker.TraceConnectStart(ctx, pgx.TraceConnectStartData{})

	// End successfully
	endData := pgx.TraceConnectEndData{Err: nil}
	tracker.TraceConnectEnd(ctxWithMeta, endData)

	// Should not panic
}

func TestPGXTracker_TraceConnectEndWithTimeout(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Start connect first
	ctxWithMeta := tracker.TraceConnectStart(ctx, pgx.TraceConnectStartData{})

	// End with timeout error
	timeoutErr := errors.New("connection timeout")
	endData := pgx.TraceConnectEndData{Err: timeoutErr}
	tracker.TraceConnectEnd(ctxWithMeta, endData)

	// Should not panic
}

func TestPGXTracker_TracePrepareStart(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	data := pgx.TracePrepareStartData{
		Name: "stmt1",
		SQL:  "SELECT $1",
	}

	newCtx := tracker.TracePrepareStart(ctx, nil, data)
	assert.NotEqual(t, ctx, newCtx)
}

func TestPGXTracker_TracePrepareEnd(t *testing.T) {
	ctx := context.Background()
	tracker := createTestPGXTracker(t)

	// Start prepare first
	startData := pgx.TracePrepareStartData{Name: "stmt1", SQL: "SELECT $1"}
	ctxWithMeta := tracker.TracePrepareStart(ctx, nil, startData)

	// End prepare
	endData := pgx.TracePrepareEndData{Err: nil}
	tracker.TracePrepareEnd(ctxWithMeta, nil, endData)

	// Should not panic
}

func TestPGXTracker_FinishMeasuringWithoutMetadata(t *testing.T) {
	tracker := createTestPGXTracker(t)
	ctx := context.Background()

	// Call finishMeasuring without proper metadata in context
	result := tracker.finishMeasuring(ctx, 0, nil)

	// Should return empty metadata
	assert.Empty(t, result.attrs)
	assert.Zero(t, result.start)
}

func TestExtractSQLOperationName(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"simple select", "SELECT * FROM users", "SELECT"},
		{"insert query", "INSERT INTO users (name) VALUES ($1)", "INSERT"},
		{"update query", "UPDATE users SET name = $1 WHERE id = $2", "UPDATE"},
		{"delete query", "DELETE FROM users WHERE id = $1", "DELETE"},
		{"with whitespace", "  \t\n  SELECT * FROM users", "SELECT"},
		{"empty query", "", ""},
		{"with parentheses", "EXECUTE(statement)", "EXECUTE"},
		{"only whitespace", "   \t\n  ", ""},
		{"single character", "X", "X"},
		{"whitespace only at start", "\n\t SELECT", "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSQLOperationName(tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractSQLCollectionName(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"select from table", "SELECT * FROM users", "users"},
		{"insert into table", "INSERT INTO products (name) VALUES ($1)", "products"},
		{"update table", "UPDATE orders SET status = $1", "orders"},
		{"with schema", "SELECT * FROM public.users", "public.users"},
		{"complex query", "SELECT u.* FROM users u JOIN orders o ON u.id = o.user_id", "users"},
		{"short query", "SELECT 1", ""},
		{"no table", "SELECT NOW()", ""},
		{"very long table name", "UPDATE " + "very_long_table_name_that_exceeds_fifty_characters_and_should_be_ignored_in_extraction" + " SET x=1", ""},
		{"table with whitespace", "SELECT * FROM   users   WHERE id = 1", "users"},
		{"table with tab", "INSERT INTO\tproducts\t(name) VALUES ($1)", ""}, // tabs not handled in extraction
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSQLCollectionName(tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractErrTypeAndResultCode(t *testing.T) {
	// Test with PostgreSQL error
	pgErr := &pgconn.PgError{
		Code: pgerrcode.UniqueViolation,
	}
	errType, resultCode := extractErrTypeAndResultCode(pgErr)
	assert.Equal(t, pgerrcode.UniqueViolation, errType)
	assert.Equal(t, pgerrcode.UniqueViolation, resultCode)

	// Test with generic error
	genericErr := errors.New("some error")
	errType, resultCode = extractErrTypeAndResultCode(genericErr)
	assert.Equal(t, "_OTHER", errType) // dbconv.ErrorTypeOther is "_OTHER"
	assert.Equal(t, "ERROR", resultCode)
}

func TestBuildQueryMetadata(t *testing.T) {
	tests := []struct {
		name     string
		opName   string
		query    string
		expected pgQueryMetadata
	}{
		{
			name:   "with operation and query",
			opName: "SELECT",
			query:  "SELECT * FROM users",
			expected: pgQueryMetadata{
				spanName: "SELECT users",
				attrs: []attribute.KeyValue{
					semconv.DBQueryText("SELECT * FROM users"),
					semconv.DBCollectionName("users"),
					semconv.DBOperationName("SELECT"),
					semconv.DBQuerySummary("SELECT users"),
				},
			},
		},
		{
			name:   "query only",
			opName: "",
			query:  "INSERT INTO products (name) VALUES ($1)",
			expected: pgQueryMetadata{
				spanName: "INSERT products",
				attrs: []attribute.KeyValue{
					semconv.DBQueryText("INSERT INTO products (name) VALUES ($1)"),
					semconv.DBCollectionName("products"),
					semconv.DBOperationName("INSERT"),
					semconv.DBQuerySummary("INSERT products"),
				},
			},
		},
		{
			name:   "operation only",
			opName: "BATCH",
			query:  "",
			expected: pgQueryMetadata{
				spanName: "BATCH",
				attrs: []attribute.KeyValue{
					semconv.DBOperationName("BATCH"),
					semconv.DBQuerySummary("BATCH"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildQueryMetadata(tt.opName, tt.query)

			assert.Equal(t, tt.expected.spanName, result.spanName)
			assert.Len(t, result.attrs, len(tt.expected.attrs))

			// Check that all expected attributes are present
			for _, expectedAttr := range tt.expected.attrs {
				found := false
				for _, actualAttr := range result.attrs {
					if actualAttr.Key == expectedAttr.Key && actualAttr.Value.AsString() == expectedAttr.Value.AsString() {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected attribute %s not found", expectedAttr.Key)
			}

			assert.True(t, time.Since(result.start) < time.Second)
		})
	}
}

func TestIsConnectionTimeoutError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"timeout error", errors.New("connection timeout"), true},
		{"context deadline", errors.New("context deadline exceeded"), true},
		{"connection refused", errors.New("connection refused"), true},
		{"generic error", errors.New("some other error"), false},
		{"uppercase timeout", errors.New("CONNECTION TIMEOUT"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isConnectionTimeoutError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create test PGX tracker using environment variables
func createTestPGXTracker(t *testing.T) *PGXTracker {
	ctx := context.Background()
	host := getEnvOrDefault("PG_HOST", "localhost")
	port := getEnvInt("PG_PORT", 5432)
	database := getEnvOrDefault("PG_DATABASE", "testdb")
	tracker, err := NewPGXTracker(ctx, host, port, database)
	require.NoError(t, err)
	return tracker
}

// Helper functions for environment variables (duplicated from integration_test.go)
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func TestFinishMeasuring_OverflowProtection(t *testing.T) {
	tests := []struct {
		name         string
		rowsAffected int64
		expectedRows int
	}{
		{
			name:         "normal rows affected should use actual value",
			rowsAffected: 1000,
			expectedRows: 1000,
		},
		{
			name:         "max int64 should be capped to max int",
			rowsAffected: 9223372036854775807, // max int64
			expectedRows: int(^uint(0) >> 1),  // max int
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Set up metadata in context
			meta := pgQueryMetadata{
				attrs:  []attribute.KeyValue{},
				start:  time.Now().Add(-100 * time.Millisecond),
				opName: "SELECT", // This will not be in opsNotReturningRows
			}
			ctx = context.WithValue(ctx, pgQueryMetadataKey{}, meta)

			// Create a mock tracer
			tracker := &PGXTracker{
				tracer: noop.NewTracerProvider().Tracer("test"),
			}

			// Call finishMeasuring
			result := tracker.finishMeasuring(ctx, tt.rowsAffected, nil)

			// Verify the overflow protection worked
			assert.NotNil(t, result.rowsAffected)
			assert.Equal(t, tt.rowsAffected, *result.rowsAffected)

			// Check that the attributes contain the expected capped value
			var foundRows bool
			for _, attr := range result.attrs {
				if attr.Key == "db.response.returned_rows" {
					foundRows = true
					assert.Equal(t, int64(tt.expectedRows), attr.Value.AsInt64())
					break
				}
			}
			assert.True(t, foundRows, "should have db.response.returned_rows attribute")
		})
	}
}
