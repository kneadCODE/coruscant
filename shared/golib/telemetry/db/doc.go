// Package db provides database observability instrumentation for PostgreSQL.
//
// This package implements automatic telemetry collection using pgx tracer hooks,
// providing comprehensive metrics and distributed tracing without manual instrumentation.
//
// Key Features:
//   - Zero-configuration observability via PGXTracker
//   - Full OpenTelemetry semantic convention compliance
//   - Perfect timing alignment between spans and metrics
//   - Automatic SQL operation and table name extraction
//   - PostgreSQL-specific error code handling
//   - Complete connection pool lifecycle monitoring
//
// Usage:
//
//	tracker, err := db.NewPGXTracker(ctx, "localhost", 5432, "mydb")
//	poolConfig.ConnConfig.Tracer = tracker
package db
