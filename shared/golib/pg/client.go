package pg

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
	dbtelemetry "github.com/kneadCODE/coruscant/shared/golib/telemetry/db"
)

// Client provides PostgreSQL database operations with telemetry and retry support
type Client struct {
	pool             *pgxpool.Pool
	options          *options
	tracker          *dbtelemetry.PGXTracker
	backoff          backoff.BackOff
	maxRetryAttempts int
}

// NewClient creates a new PostgreSQL client with the given options
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	if err := options.validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// Initialize metrics BEFORE pool creation for tracer setup
	tracker, err := dbtelemetry.NewPGXTracker(ctx, options.host, options.port, options.database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database metrics: %w", err)
	}

	// Create backoff strategy
	var bo backoff.BackOff
	var maxAttempts int
	if options.enableRetry {
		bo, maxAttempts = newBackoff(options.retryDelay, options.maxRetryDelay, options.maxRetryAttempts)
	} else {
		bo = &backoff.StopBackOff{} // No retry by default
		maxAttempts = 1
	}

	pool, err := createPool(ctx, tracker, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	client := &Client{
		pool:             pool,
		options:          options,
		tracker:          tracker,
		backoff:          bo,
		maxRetryAttempts: maxAttempts,
	}

	// Record initial pool configuration metrics
	client.RecordPoolMetrics(ctx)

	// Log successful client initialization
	telemetry.RecordInfoEvent(ctx, "PostgreSQL client initialized successfully",
		"host", options.host,
		"port", options.port,
		"database", options.database,
		"max_conns", options.maxConns,
		"min_conns", options.minConns,
		"retry_enabled", options.enableRetry,
	)

	return client, nil
}

// Close closes the database connection pool
func (c *Client) Close() {
	c.pool.Close()
}

// RecordPoolMetrics records current connection pool metrics
// This method should be called periodically to maintain current pool state visibility
func (c *Client) RecordPoolMetrics(ctx context.Context) {
	if c.tracker == nil {
		return
	}

	stats := c.pool.Stat()
	if stats == nil {
		return
	}

	// Record connection counts by state (current values)
	c.tracker.RecordConnectionStats(ctx, "idle", int64(stats.IdleConns()))
	c.tracker.RecordConnectionStats(ctx, "used", int64(stats.AcquiredConns()))

	// Record pool configuration metrics (current limits)
	c.tracker.SetConnectionMax(ctx, int64(stats.MaxConns()))
	c.tracker.SetConnectionIdleMax(ctx, int64(c.options.maxConns)) // Use configured max as idle max
	c.tracker.SetConnectionIdleMin(ctx, int64(c.options.minConns)) // Use configured min as idle min

	// Record current pending requests
	c.tracker.RecordConnectionPendingRequests(ctx, stats.EmptyAcquireCount())
}

// Ping verifies a connection to the database is still alive
func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

// Query executes a query that returns rows with timeout
func (c *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	// Apply query timeout if not already set
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && c.options.queryTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.queryTimeout)
		defer cancel()
	}

	var rows pgx.Rows
	err := c.retryOperation(ctx, retryableRead, func() error {
		var err error
		rows, err = c.pool.Query(ctx, sql, args...)
		return err
	})
	return rows, err
}

// QueryRow executes a query that returns at most one row with timeout
func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	// Apply query timeout if not already set
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && c.options.queryTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.queryTimeout)
		defer cancel()
	}

	// For QueryRow, we can't use retry logic because we need to return the Row immediately
	return c.pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query without returning any rows with timeout
// Note: Exec operations are NOT retried automatically to prevent duplicates
func (c *Client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	// Apply query timeout if not already set
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && c.options.queryTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.options.queryTimeout)
		defer cancel()
	}

	return c.pool.Exec(ctx, sql, args...)
}
