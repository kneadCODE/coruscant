package pg

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// RetryableOperations defines which operations can be safely retried
type retryableOperations int

const (
	retryableRead  retryableOperations = iota // SELECT operations
	retryableWrite                            // INSERT/UPDATE/DELETE operations (dangerous!)
)

// retryOperation executes an operation with retry logic if enabled
func (c *Client) retryOperation(ctx context.Context, opType retryableOperations, operation func() error) error {
	if !c.options.enableRetry {
		return operation()
	}

	return retryOperation(ctx, c.backoff, c.maxRetryAttempts, opType, operation)
}

// errorRow is a pgx.Row implementation that always returns an error
type errorRow struct {
	err error
}

// Scan returns the stored error, implementing pgx.Row interface
func (r *errorRow) Scan(dest ...any) error {
	return r.err
}

// IsRetryableError determines if a PostgreSQL error should be retried
func isRetryableError(err error, opType retryableOperations) bool {
	if err == nil {
		return false
	}

	// Check for context cancellation - never retry
	if errors.Is(err, context.Canceled) {
		return false
	}

	// Handle pgx specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return isPostgreSQLErrorRetryable(pgErr.Code, opType)
	}

	// Network timeouts - only retry reads to avoid duplicates
	if errors.Is(err, context.DeadlineExceeded) {
		return opType == retryableRead
	}

	return false
}

// isPostgreSQLErrorRetryable determines if a PostgreSQL error code should be retried
func isPostgreSQLErrorRetryable(code string, opType retryableOperations) bool {
	// Serialization and concurrency errors - safe to retry for all operations
	if isSerializationError(code) {
		return true
	}

	// Connection and resource issues - only retry reads to avoid duplicate writes
	if isConnectionError(code) || isResourceError(code) || isAdminError(code) {
		return opType == retryableRead
	}

	return false
}

// isSerializationError checks if error is a serialization/concurrency issue
func isSerializationError(code string) bool {
	return code == pgerrcode.SerializationFailure || code == pgerrcode.DeadlockDetected
}

// isConnectionError checks if error is related to connection issues
func isConnectionError(code string) bool {
	return code == pgerrcode.ConnectionException ||
		code == pgerrcode.ConnectionDoesNotExist ||
		code == pgerrcode.ConnectionFailure ||
		code == pgerrcode.SQLClientUnableToEstablishSQLConnection ||
		code == pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection
}

// isResourceError checks if error is related to resource exhaustion
func isResourceError(code string) bool {
	return code == pgerrcode.InsufficientResources ||
		code == pgerrcode.DiskFull ||
		code == pgerrcode.OutOfMemory ||
		code == pgerrcode.TooManyConnections
}

// isAdminError checks if error is related to admin operations
func isAdminError(code string) bool {
	return code == pgerrcode.AdminShutdown
}

// NewBackoff creates a backoff strategy for retries
func newBackoff(initialDelay, maxDelay time.Duration, maxAttempts int) (backoff.BackOff, int) {
	// Validate maxAttempts to ensure safe conversion
	if maxAttempts <= 0 {
		maxAttempts = 1 // Default to single attempt for safety
	}

	exponential := backoff.NewExponentialBackOff()
	exponential.InitialInterval = initialDelay
	exponential.MaxInterval = maxDelay
	exponential.Multiplier = 2.0
	exponential.RandomizationFactor = 0.1 // 10% jitter

	return exponential, maxAttempts
}

// RetryOperation executes an operation with retry logic using cenkalti/backoff
func retryOperation(ctx context.Context, bo backoff.BackOff, maxAttempts int, opType retryableOperations, operation func() error) error {
	attempt := 0
	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		attempt++
		err := operation()

		if err == nil {
			if attempt > 1 {
				telemetry.RecordInfoEvent(ctx, "Database operation succeeded after retry",
					"attempts", attempt,
					"operation_type", opTypeString(opType),
				)
			}
			return struct{}{}, nil
		}

		// Check if error is retryable for this operation type
		if !isRetryableError(err, opType) {
			telemetry.RecordDebugEvent(ctx, "Database operation failed with non-retryable error",
				"error", err.Error(),
				"operation_type", opTypeString(opType),
				"attempts", attempt,
			)
			return struct{}{}, backoff.Permanent(err) // Don't retry non-retryable errors
		}

		telemetry.RecordDebugEvent(ctx, "Database operation failed, will retry",
			"error", err.Error(),
			"operation_type", opTypeString(opType),
			"attempt", attempt,
		)

		return struct{}{}, err
	}, backoff.WithBackOff(bo), backoff.WithMaxTries(uint(maxAttempts))) // #nosec G115 - Safe after explicit validation
	return err
}

// opTypeString converts operation type to string for logging
func opTypeString(opType retryableOperations) string {
	switch opType {
	case retryableRead:
		return "read"
	case retryableWrite:
		return "write"
	default:
		return "unknown"
	}
}
