package pg

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
	dbtelemetry "github.com/kneadCODE/coruscant/shared/golib/telemetry/db"
)

func createPool(ctx context.Context, tracker *dbtelemetry.PGXTracker, opts *options) (*pgxpool.Pool, error) {
	dsn := buildDSN(opts)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = int32(opts.maxConns)
	poolConfig.MinConns = int32(opts.minConns)
	poolConfig.MaxConnLifetime = opts.maxConnLifetime
	poolConfig.MaxConnIdleTime = opts.maxConnIdleTime
	poolConfig.ConnConfig.Tracer = tracker

	poolConfig.BeforeConnect = func(ctx context.Context, config *pgx.ConnConfig) error {
		// Connection lifecycle is already tracked by TraceConnectStart/End
		return nil
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Connection creation metrics are handled by TraceConnectEnd
		telemetry.RecordDebugEvent(ctx, "Database connection established successfully")
		return nil
	}
	poolConfig.BeforeClose = func(conn *pgx.Conn) {
		// Connection close metrics are tracked through pool statistics
		telemetry.RecordDebugEvent(ctx, "Closing database connection")
	}
	poolConfig.AfterRelease = func(conn *pgx.Conn) bool {
		// Return connection health status to pool manager
		isHealthy := !conn.IsClosed()
		if !isHealthy {
			telemetry.RecordDebugEvent(context.Background(), "Discarding unhealthy connection from pool")
		}
		return isHealthy
	}

	poolCreateStart := time.Now()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	poolCreateDuration := time.Since(poolCreateStart)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	tracker.RecordConnectionCreateTime(ctx, poolCreateDuration)

	return pool, nil
}

func buildDSN(opts *options) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(opts.username, opts.password),
		Host:   fmt.Sprintf("%s:%d", opts.host, opts.port),
		Path:   opts.database,
	}

	query := url.Values{}
	query.Set("sslmode", opts.sslMode)
	query.Set("connect_timeout", fmt.Sprintf("%.0f", opts.connectTimeout.Seconds()))
	u.RawQuery = query.Encode()

	return u.String()
}
