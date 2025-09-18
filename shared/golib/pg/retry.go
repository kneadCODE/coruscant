package pg

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
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

	return retryOperation(ctx, c.backoff, opType, operation)
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
		switch pgErr.Code {
		// Serialization and concurrency errors - safe to retry for all operations
		case pgerrcode.SerializationFailure:
			return true
		case pgerrcode.DeadlockDetected:
			return true

		// Connection issues - safe to retry for reads, risky for writes
		case pgerrcode.ConnectionException,
			pgerrcode.ConnectionDoesNotExist,
			pgerrcode.ConnectionFailure,
			pgerrcode.SQLClientUnableToEstablishSQLConnection,
			pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection:
			return opType == retryableRead

		// Resource exhaustion - only retry reads to avoid duplicate writes
		case pgerrcode.InsufficientResources,
			pgerrcode.DiskFull,
			pgerrcode.OutOfMemory,
			pgerrcode.TooManyConnections:
			return opType == retryableRead

		// Admin operations - only retry reads
		case pgerrcode.AdminShutdown:
			return opType == retryableRead
		}
	}

	// Network timeouts - only retry reads to avoid duplicates
	if errors.Is(err, context.DeadlineExceeded) {
		return opType == retryableRead
	}

	return false
}

// NewBackoff creates a backoff strategy for retries
func newBackoff(initialDelay, maxDelay time.Duration, maxAttempts int) backoff.BackOff {
	exponential := backoff.NewExponentialBackOff()
	exponential.InitialInterval = initialDelay
	exponential.MaxInterval = maxDelay
	exponential.MaxElapsedTime = 0 // Disable time-based limit, use attempt-based only
	exponential.Multiplier = 2.0
	exponential.RandomizationFactor = 0.1 // 10% jitter

	return backoff.WithMaxRetries(exponential, uint64(maxAttempts-1)) // #nosec G115 - maxAttempts is validated positive
}

// RetryOperation executes an operation with retry logic using cenkalti/backoff
func retryOperation(ctx context.Context, bo backoff.BackOff, opType retryableOperations, operation func() error) error {
	attempt := 0
	return backoff.Retry(func() error {
		attempt++
		err := operation()

		if err == nil {
			if attempt > 1 {
				telemetry.RecordInfoEvent(ctx, "Database operation succeeded after retry",
					"attempts", attempt,
					"operation_type", opTypeString(opType),
				)
			}
			return nil
		}

		// Check if error is retryable for this operation type
		if !isRetryableError(err, opType) {
			telemetry.RecordDebugEvent(ctx, "Database operation failed with non-retryable error",
				"error", err.Error(),
				"operation_type", opTypeString(opType),
				"attempts", attempt,
			)
			return backoff.Permanent(err) // Don't retry non-retryable errors
		}

		telemetry.RecordDebugEvent(ctx, "Database operation failed, will retry",
			"error", err.Error(),
			"operation_type", opTypeString(opType),
			"attempt", attempt,
		)

		return err
	}, backoff.WithContext(bo, ctx))
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
