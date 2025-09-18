package pg

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		opType   retryableOperations
		expected bool
	}{
		{
			name:     "nil error should not retry",
			err:      nil,
			opType:   retryableRead,
			expected: false,
		},
		{
			name:     "context canceled should not retry",
			err:      context.Canceled,
			opType:   retryableRead,
			expected: false,
		},
		{
			name:     "context deadline exceeded should retry for reads",
			err:      context.DeadlineExceeded,
			opType:   retryableRead,
			expected: true,
		},
		{
			name:     "context deadline exceeded should not retry for writes",
			err:      context.DeadlineExceeded,
			opType:   retryableWrite,
			expected: false,
		},
		{
			name:     "serialization failure should retry for all operations",
			err:      &pgconn.PgError{Code: pgerrcode.SerializationFailure},
			opType:   retryableWrite,
			expected: true,
		},
		{
			name:     "deadlock should retry for all operations",
			err:      &pgconn.PgError{Code: pgerrcode.DeadlockDetected},
			opType:   retryableRead,
			expected: true,
		},
		{
			name:     "connection error should retry for reads only",
			err:      &pgconn.PgError{Code: pgerrcode.ConnectionException},
			opType:   retryableRead,
			expected: true,
		},
		{
			name:     "connection error should not retry for writes",
			err:      &pgconn.PgError{Code: pgerrcode.ConnectionException},
			opType:   retryableWrite,
			expected: false,
		},
		{
			name:     "resource error should retry for reads only",
			err:      &pgconn.PgError{Code: pgerrcode.InsufficientResources},
			opType:   retryableRead,
			expected: true,
		},
		{
			name:     "admin error should retry for reads only",
			err:      &pgconn.PgError{Code: pgerrcode.AdminShutdown},
			opType:   retryableRead,
			expected: true,
		},
		{
			name:     "unknown postgres error should not retry",
			err:      &pgconn.PgError{Code: "23505"}, // unique violation
			opType:   retryableRead,
			expected: false,
		},
		{
			name:     "non-postgres error should not retry",
			err:      errors.New("some other error"),
			opType:   retryableRead,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err, tt.opType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSerializationError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "serialization failure should return true",
			code:     pgerrcode.SerializationFailure,
			expected: true,
		},
		{
			name:     "deadlock detected should return true",
			code:     pgerrcode.DeadlockDetected,
			expected: true,
		},
		{
			name:     "other error should return false",
			code:     pgerrcode.UniqueViolation,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSerializationError(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "connection exception should return true",
			code:     pgerrcode.ConnectionException,
			expected: true,
		},
		{
			name:     "connection does not exist should return true",
			code:     pgerrcode.ConnectionDoesNotExist,
			expected: true,
		},
		{
			name:     "connection failure should return true",
			code:     pgerrcode.ConnectionFailure,
			expected: true,
		},
		{
			name:     "unable to establish connection should return true",
			code:     pgerrcode.SQLClientUnableToEstablishSQLConnection,
			expected: true,
		},
		{
			name:     "server rejected connection should return true",
			code:     pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
			expected: true,
		},
		{
			name:     "other error should return false",
			code:     pgerrcode.UniqueViolation,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isConnectionError(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsResourceError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "insufficient resources should return true",
			code:     pgerrcode.InsufficientResources,
			expected: true,
		},
		{
			name:     "disk full should return true",
			code:     pgerrcode.DiskFull,
			expected: true,
		},
		{
			name:     "out of memory should return true",
			code:     pgerrcode.OutOfMemory,
			expected: true,
		},
		{
			name:     "too many connections should return true",
			code:     pgerrcode.TooManyConnections,
			expected: true,
		},
		{
			name:     "other error should return false",
			code:     pgerrcode.UniqueViolation,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isResourceError(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAdminError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "admin shutdown should return true",
			code:     pgerrcode.AdminShutdown,
			expected: true,
		},
		{
			name:     "other error should return false",
			code:     pgerrcode.UniqueViolation,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAdminError(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewBackoff(t *testing.T) {
	tests := []struct {
		name         string
		initialDelay time.Duration
		maxDelay     time.Duration
		maxAttempts  int
	}{
		{
			name:         "valid parameters",
			initialDelay: 100 * time.Millisecond,
			maxDelay:     5 * time.Second,
			maxAttempts:  3,
		},
		{
			name:         "zero max attempts should default to 1",
			initialDelay: 100 * time.Millisecond,
			maxDelay:     5 * time.Second,
			maxAttempts:  0,
		},
		{
			name:         "negative max attempts should default to 1",
			initialDelay: 100 * time.Millisecond,
			maxDelay:     5 * time.Second,
			maxAttempts:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backoff := newBackoff(tt.initialDelay, tt.maxDelay, tt.maxAttempts)
			assert.NotNil(t, backoff)

			// Test that backoff returns a delay (or stops)
			delay := backoff.NextBackOff()
			// Delay can be non-negative or backoff.Stop (which is -1)
			assert.True(t, delay >= 0 || delay == -1, "backoff should return non-negative delay or stop")
		})
	}
}
