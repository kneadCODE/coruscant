package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// TxFunc is a function that executes within a transaction context
type TxFunc func(ctx context.Context, tx pgx.Tx) error

// WithTx executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// If the function completes successfully, the transaction is committed
func (c *Client) WithTx(ctx context.Context, fn TxFunc) error {
	return c.WithTxOptions(ctx, pgx.TxOptions{}, fn)
}

// WithTxOptions executes a function within a transaction with specific options
func (c *Client) WithTxOptions(ctx context.Context, txOptions pgx.TxOptions, fn TxFunc) error {
	tx, err := c.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	fnErr := fn(ctx, tx)

	if fnErr != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			telemetry.RecordErrorEvent(ctx, fmt.Errorf("transaction rollback failed after function error: %w", rollbackErr))
		}
		return fnErr
	}

	// Function succeeded - commit transaction
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}
