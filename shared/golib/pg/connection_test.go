package pg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbtelemetry "github.com/kneadCODE/coruscant/shared/golib/telemetry/db"
)

func TestCreatePool_ValidationEdgeCases(t *testing.T) {
	ctx := context.Background()

	// Create a basic tracker for testing using environment variables
	host := getEnvOrDefault("PG_HOST", "localhost")
	port := getEnvInt("PG_PORT", 5432)
	database := getEnvOrDefault("PG_DATABASE", "testdb")
	tracker, err := dbtelemetry.NewPGXTracker(ctx, host, port, database)
	require.NoError(t, err)

	tests := []struct {
		name        string
		opts        *options
		expectError bool
		errorMsg    string
	}{
		{
			name: "maxConns at int32 limit should work",
			opts: &options{
				host:     "localhost",
				port:     5432,
				database: "test",
				username: "user",
				password: "pass",
				sslMode:  "disable",
				maxConns: 2147483647, // max int32
				minConns: 1,
			},
			expectError: false,
		},
		{
			name: "maxConns exceeding int32 limit should fail",
			opts: &options{
				host:     "localhost",
				port:     5432,
				database: "test",
				username: "user",
				password: "pass",
				sslMode:  "disable",
				maxConns: 2147483648, // exceeds int32
				minConns: 1,
			},
			expectError: true,
			errorMsg:    "maxConns value 2147483648 exceeds int32 limit",
		},
		{
			name: "minConns exceeding int32 limit should fail",
			opts: &options{
				host:     "localhost",
				port:     5432,
				database: "test",
				username: "user",
				password: "pass",
				sslMode:  "disable",
				maxConns: 25,
				minConns: 2147483648, // exceeds int32
			},
			expectError: true,
			errorMsg:    "minConns value 2147483648 exceeds int32 limit",
		},
		{
			name: "normal values should work",
			opts: &options{
				host:     getEnvOrDefault("PG_HOST", "localhost"),
				port:     getEnvInt("PG_PORT", 5432),
				database: getEnvOrDefault("PG_DATABASE", "test"),
				username: getEnvOrDefault("PG_USERNAME", "user"),
				password: getEnvOrDefault("PG_PASSWORD", "pass"),
				sslMode:  getEnvOrDefault("PG_SSL_MODE", "disable"),
				maxConns: 25,
				minConns: 5,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := createPool(ctx, tracker, tt.opts)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else if err != nil {
				// Note: This will still fail because we're not connecting to a real DB,
				// but we're testing that our validation passes before the connection attempt
				// If there's an error, it should NOT be our validation error
				assert.NotContains(t, err.Error(), "exceeds int32 limit")
			}
		})
	}
}

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		opts     *options
		expected string
	}{
		{
			name: "basic DSN construction",
			opts: &options{
				host:           "localhost",
				port:           5432,
				database:       "testdb",
				username:       "user",
				password:       "pass",
				sslMode:        "disable",
				connectTimeout: 10000000000, // 10 seconds in nanoseconds
			},
			expected: "postgres://user:pass@localhost:5432/testdb?connect_timeout=10&sslmode=disable",
		},
		{
			name: "DSN with special characters",
			opts: &options{
				host:           "db.example.com",
				port:           3306,
				database:       "my-app",
				username:       "app_user",
				password:       "complex@pass",
				sslMode:        "require",
				connectTimeout: 5000000000, // 5 seconds
			},
			expected: "postgres://app_user:complex%40pass@db.example.com:3306/my-app?connect_timeout=5&sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildDSN(tt.opts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
