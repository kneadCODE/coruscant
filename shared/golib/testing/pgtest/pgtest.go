// Package pgtest provides PostgreSQL testing utilities with transaction-based isolation
package pgtest

import (
	"context"
	"errors"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kneadCODE/coruscant/shared/golib/pg"
)

// TB is an interface that both *testing.T and *testing.B implement
type TB interface {
	Helper()
	Skip(args ...any)
}

var (
	// Package-level connection, initialized once
	testClient *pg.Client
	clientOnce sync.Once
	clientErr  error
)

// TestFunc represents a function that runs within a test transaction
type TestFunc func(ctx context.Context, tx pgx.Tx) error

// rollbackError is a sentinel error used to force transaction rollback for test isolation
type rollbackError struct{}

func (rollbackError) Error() string { return "test isolation rollback" }

// TestWithDB executes a test function within a transaction that always rolls back
// This ensures complete test isolation - no test data persists after execution
func TestWithDB(t *testing.T, fn TestFunc) {
	t.Helper()

	client := getTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use WithTx but always force rollback by returning an error
	err := client.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Execute test function
		if err := fn(ctx, tx); err != nil {
			t.Fatalf("Test function failed: %v", err)
		}

		// Always return error to force rollback for test isolation
		return rollbackError{}
	})

	// Check if error is our intentional rollback or a real error
	var rollbackError rollbackError
	if err != nil && !errors.As(err, &rollbackError) {
		t.Fatalf("Transaction failed: %v", err)
	}
}

// getTestClient returns the shared test client, initializing it if needed
func getTestClient(t TB) *pg.Client {
	t.Helper()

	clientOnce.Do(func() {
		testClient, clientErr = createTestClient()
	})

	if clientErr != nil {
		t.Skip("PostgreSQL not available for testing:", clientErr)
	}

	return testClient
}

// createTestClient creates a new test database client
func createTestClient() (*pg.Client, error) {
	options := []pg.Option{
		pg.WithHost(getEnvOrDefault("PG_HOST", "localhost")),
		pg.WithPort(getEnvInt("PG_PORT", 5432)),
		pg.WithDatabase(getEnvOrDefault("PG_DATABASE", "postgres")),
		pg.WithCredentials(
			getEnvOrDefault("PG_USERNAME", "coruscant"),
			getEnvOrDefault("PG_PASSWORD", "trust"), // Dummy password for trust auth
		),
		pg.WithSSLMode("disable"),
		pg.WithMaxConnections(10), // Lower for tests
		pg.WithMinConnections(2),
		pg.WithConnectTimeout(5 * time.Second),
		pg.WithQueryTimeout(10 * time.Second), // Reasonable test timeout
		pg.WithoutRetry(),                     // Disable retry for predictable tests
	}

	client, err := pg.NewClient(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	// Test connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// RequireDB ensures PostgreSQL is available, skipping the test if not
func RequireDB(t *testing.T) *pg.Client {
	t.Helper()
	return getTestClient(t)
}

// RequireDBForBenchmark ensures PostgreSQL is available for benchmark tests
func RequireDBForBenchmark(b *testing.B) *pg.Client {
	b.Helper()
	return getTestClient(b)
}

// CleanupDB closes the shared test client (call in TestMain if needed)
func CleanupDB() {
	if testClient != nil {
		testClient.Close()
	}
}

// Helper functions
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
