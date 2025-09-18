package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/dbconv"
)

func TestBuildCommonAttrs(t *testing.T) {
	tests := []struct {
		name     string
		dbSystem dbconv.SystemNameAttr
		host     string
		port     int
		database string
		expected []attribute.KeyValue
	}{
		{
			name:     "PostgreSQL database",
			dbSystem: dbconv.SystemNamePostgreSQL,
			host:     "localhost",
			port:     5432,
			database: "testdb",
			expected: []attribute.KeyValue{
				attribute.String("db.system.name", "postgresql"),
				semconv.ServerAddress("localhost"),
				semconv.ServerPort(5432),
				attribute.String("db.namespace", "testdb"),
			},
		},
		{
			name:     "MySQL database",
			dbSystem: dbconv.SystemNameMySQL,
			host:     "db.example.com",
			port:     3306,
			database: "myapp",
			expected: []attribute.KeyValue{
				attribute.String("db.system.name", "mysql"),
				semconv.ServerAddress("db.example.com"),
				semconv.ServerPort(3306),
				attribute.String("db.namespace", "myapp"),
			},
		},
		{
			name:     "MongoDB database",
			dbSystem: dbconv.SystemNameMongoDB,
			host:     "mongo.internal",
			port:     27017,
			database: "analytics",
			expected: []attribute.KeyValue{
				attribute.String("db.system.name", "mongodb"),
				semconv.ServerAddress("mongo.internal"),
				semconv.ServerPort(27017),
				attribute.String("db.namespace", "analytics"),
			},
		},
		{
			name:     "Redis with default port",
			dbSystem: dbconv.SystemNameRedis,
			host:     "redis-cluster",
			port:     6379,
			database: "cache",
			expected: []attribute.KeyValue{
				attribute.String("db.system.name", "redis"),
				semconv.ServerAddress("redis-cluster"),
				semconv.ServerPort(6379),
				attribute.String("db.namespace", "cache"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := buildCommonAttrs(tt.dbSystem, tt.host, tt.port, tt.database)
			assert.Equal(t, tt.expected, attrs)
		})
	}
}

func TestInstrumentationIdentifier(t *testing.T) {
	assert.Equal(t, "github.com/kneadCODE/coruscant/shared/golib/telemetry", instrumentationIdentifier)
}
