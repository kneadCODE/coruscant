package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/dbconv"
)

func TestBuildCommonAttrs(t *testing.T) {
	attrs := buildCommonAttrs(dbconv.SystemNamePostgreSQL, "localhost", 5432, "testdb")

	expected := []attribute.KeyValue{
		attribute.String("db.system.name", "postgresql"),
		semconv.ServerAddress("localhost"),
		semconv.ServerPort(5432),
		attribute.String("db.namespace", "testdb"),
	}

	assert.Equal(t, expected, attrs)
}

func TestBuildCommonAttrsWithDifferentDB(t *testing.T) {
	attrs := buildCommonAttrs(dbconv.SystemNameMySQL, "db.example.com", 3306, "myapp")

	expected := []attribute.KeyValue{
		attribute.String("db.system.name", "mysql"),
		semconv.ServerAddress("db.example.com"),
		semconv.ServerPort(3306),
		attribute.String("db.namespace", "myapp"),
	}

	assert.Equal(t, expected, attrs)
}

func TestInstrumentationIdentifier(t *testing.T) {
	assert.Equal(t, "github.com/kneadCODE/coruscant/shared/golib/telemetry", instrumentationIdentifier)
}
