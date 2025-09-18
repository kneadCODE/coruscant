package pg

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_OptionsPattern(t *testing.T) {
	testOptions()
	client, err := NewClient(context.Background(),
		WithHost(os.Getenv("PG_HOST")),
		WithPort(5432),
		WithDatabase(os.Getenv("PG_DATABASE")),
		WithCredentials(os.Getenv("PG_USERNAME"), "trust"),
		WithSSLMode("disable"),
		WithMaxConnections(5),
		WithMinConnections(1),
		WithConnectTimeout(5*time.Second),
		WithoutRetry(),
	)
	if err != nil {
		t.Skip("PostgreSQL not available, skipping test")
		return
	}
	defer client.Close()

	// Test basic functionality
	ctx := context.Background()
	err = client.Ping(ctx)
	assert.NoError(t, err)
}

func TestClient_BasicOperations(t *testing.T) {
	skipIfNoPostgreSQL(t)

	client := setupTestDB(t)
	defer cleanupTestDB(t, client)

	ctx := context.Background()

	// Test ping
	err := client.Ping(ctx)
	assert.NoError(t, err)

	// Test simple query
	var version string
	err = client.QueryRow(ctx, "SELECT version()").Scan(&version)
	assert.NoError(t, err)
	assert.Contains(t, version, "PostgreSQL")

	// Test transaction
	err = client.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var result int
		return tx.QueryRow(ctx, "SELECT 42").Scan(&result)
	})
	assert.NoError(t, err)
}

func TestClient_TransactionRollback(t *testing.T) {
	skipIfNoPostgreSQL(t)

	client := setupTestDB(t)
	defer cleanupTestDB(t, client)

	ctx := context.Background()

	// Create a simple temporary table for testing
	_, err := client.Exec(ctx, `
		CREATE TEMP TABLE test_rollback (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	// Transaction that should rollback
	err = client.WithTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			"INSERT INTO test_rollback (name) VALUES ($1)",
			"Test User")
		if err != nil {
			return err
		}

		// Return error to trigger rollback
		return assert.AnError
	})
	assert.Error(t, err)

	// Verify no data was inserted (rollback worked)
	var count int
	err = client.QueryRow(ctx, "SELECT COUNT(*) FROM test_rollback").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestClient_Ping(t *testing.T) {
	skipIfNoPostgreSQL(t)

	client := setupTestDB(t)
	defer cleanupTestDB(t, client)

	// Test ping
	err := client.Ping(context.Background())
	assert.NoError(t, err)

	// Test ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx)
	assert.NoError(t, err)
}

func TestClient_InvalidOptions(t *testing.T) {
	// Test missing required options
	_, err := NewClient(context.Background(),
		WithHost("localhost"),
		// Missing database, username, password
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid options")
}

// testOptions returns test options for PostgreSQL
func testOptions() []Option {
	// Try to get config from environment first
	host := os.Getenv("PG_HOST")
	port := getEnvInt("PG_PORT", 5432)
	database := os.Getenv("PG_DATABASE")
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD") // Dummy password for trust auth

	return []Option{
		WithHost(host),
		WithPort(port),
		WithDatabase(database),
		WithCredentials(username, password),
		WithSSLMode("disable"),
		WithMaxConnections(10), // Lower for tests
		WithMinConnections(2),
		WithConnectTimeout(5 * time.Second),
		WithoutRetry(), // Disable retry for predictable tests
	}
}

// setupTestDB creates a test database client and ensures connectivity
func setupTestDB(t *testing.T) *Client {
	t.Helper()

	client, err := NewClient(context.Background(), testOptions()...)
	require.NoError(t, err, "Failed to create test database client")

	// Test connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx)
	require.NoError(t, err, "Failed to ping test database")

	return client
}

// cleanupTestDB closes the database connection
func cleanupTestDB(t *testing.T, client *Client) {
	t.Helper()
	if client != nil {
		client.Close()
	}
}

// skipIfNoPostgreSQL skips the test if PostgreSQL is not available
func skipIfNoPostgreSQL(t *testing.T) {
	t.Helper()

	client, err := NewClient(context.Background(), testOptions()...)
	if err != nil {
		t.Skip("PostgreSQL not available, skipping test")
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		t.Skip("PostgreSQL not available, skipping test")
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
