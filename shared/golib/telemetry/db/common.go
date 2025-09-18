package db

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/dbconv"
)

const (
	// instrumentationIdentifier is used for OpenTelemetry instrumentation naming
	instrumentationIdentifier = "github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// buildCommonAttrs creates the base set of attributes applied to all database telemetry.
// These attributes satisfy OpenTelemetry semantic convention requirements for database operations.
func buildCommonAttrs(
	dbSystem dbconv.SystemNameAttr, serverAddr string, serverPort int, dbName string,
) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.DBSystemNameKey.String(string(dbSystem)),
		semconv.ServerAddress(serverAddr),
		semconv.ServerPort(serverPort),
		semconv.DBNamespace(dbName),
	}
}
