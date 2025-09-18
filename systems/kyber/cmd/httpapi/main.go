package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kneadCODE/coruscant/shared/golib/httpserver"
	"github.com/kneadCODE/coruscant/shared/golib/pg"
	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// Global database client (initialized in main)
var dbClient *pg.Client

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Initialize telemetry with dev debug mode for comprehensive observability
	ctx, cleanup, err := telemetry.InitTelemetry(ctx, telemetry.ModeDebug)
	if err != nil {
		return err
	}
	defer cleanup(context.Background())

	// Initialize database connection with automatic observability
	dbClient, err = pg.NewClient(ctx,
		pg.WithHost(os.Getenv("PG_HOST")),
		pg.WithPort(getEnvInt("PG_PORT")),
		pg.WithDatabase(os.Getenv("PG_DATABASE")),
		pg.WithCredentials(
			os.Getenv("PG_USERNAME"),
			os.Getenv("PG_PASSWORD"),
		),
		pg.WithSSLMode(os.Getenv("PG_SSL_MODE")),
		pg.WithMaxConnections(10),
		pg.WithMinConnections(2),
		pg.WithRetrySettings(3, 100*time.Millisecond, 5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create database client: %w", err)
	}
	defer dbClient.Close()

	// Test database connection
	if err := dbClient.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create database schema (in production, use migrations)
	if err := createSchema(ctx, dbClient); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	telemetry.RecordInfoEvent(ctx, "Database initialized successfully")

	return start(ctx)
}

func start(ctx context.Context) error {
	// Create HTTP server with comprehensive observability features enabled
	srv, err := httpserver.NewServer(ctx,
		httpserver.WithRESTHandler(restHandler),
		httpserver.WithProfilingHandler(), // Enable pprof profiling endpoints
		httpserver.WithMetricsHandler(),   // Enable metrics endpoint (shows OTEL info)
	)
	if err != nil {
		return err
	}
	return srv.Start(ctx)
}

func restHandler(rtr chi.Router) {
	userService := NewUserService(dbClient)
	setupRoutes(rtr, userService)
}

func getEnvInt(key string) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return 0
}
