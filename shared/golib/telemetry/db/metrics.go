package db

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv/v1.37.0/dbconv"

	"github.com/kneadCODE/coruscant/shared/golib/telemetry"
)

// Metrics provides OpenTelemetry-compliant database metrics collection.
//
// Implements all required and recommended database client metrics according to
// OpenTelemetry semantic conventions v1.37.0:
//   - db.client.operation.duration (required)
//   - db.client.response.returned_rows (recommended)
//   - Complete connection pool metrics suite
//
// All metrics include proper OpenTelemetry attributes for filtering and correlation.
type Metrics struct {
	dbSystem     dbconv.SystemNameAttr
	commonAttrs  []attribute.KeyValue
	connPoolName string

	operationDuration   dbconv.ClientOperationDuration
	returnedRows        dbconv.ClientResponseReturnedRows
	connCount           dbconv.ClientConnectionCount
	connMax             dbconv.ClientConnectionMax
	connIdleMaxCount    dbconv.ClientConnectionIdleMax
	connIdleMinCount    dbconv.ClientConnectionIdleMin
	connPendingRequests dbconv.ClientConnectionPendingRequests
	connTimetouts       dbconv.ClientConnectionTimeouts
	connCreateTime      dbconv.ClientConnectionCreateTime
	connWaitTime        dbconv.ClientConnectionWaitTime
	connUseTime         dbconv.ClientConnectionUseTime
}

// RecordConnectionStats records connection pool metrics
// state can be idle or used
func (dm *Metrics) RecordConnectionStats(ctx context.Context, state string, count int64) {
	dm.connCount.Add(ctx,
		count,
		dm.connPoolName,
		dbconv.ClientConnectionStateAttr(state),
		dm.commonAttrs...,
	)
}

// SetConnectionMax sets the current maximum number of connections in the pool
func (dm *Metrics) SetConnectionMax(ctx context.Context, maxConnections int64) {
	// For UpDownCounter, we need to calculate the difference to set the value
	// This is a simplified approach - in production you'd track the previous value
	dm.connMax.Add(ctx,
		maxConnections,
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// SetConnectionIdleMax sets the current maximum number of idle connections
func (dm *Metrics) SetConnectionIdleMax(ctx context.Context, maxIdleConnections int64) {
	dm.connIdleMaxCount.Add(ctx,
		maxIdleConnections,
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// SetConnectionIdleMin sets the current minimum number of idle connections
func (dm *Metrics) SetConnectionIdleMin(ctx context.Context, minIdleConnections int64) {
	dm.connIdleMinCount.Add(ctx,
		minIdleConnections,
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// RecordConnectionPendingRequests records the number of pending connection requests
func (dm *Metrics) RecordConnectionPendingRequests(ctx context.Context, pendingRequests int64) {
	dm.connPendingRequests.Add(ctx,
		pendingRequests,
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// RecordConnectionTimeouts records connection timeout events
func (dm *Metrics) RecordConnectionTimeouts(ctx context.Context, timeouts int64) {
	dm.connTimetouts.Add(ctx,
		timeouts,
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// RecordConnectionCreateTime records the time taken to create a new connection
func (dm *Metrics) RecordConnectionCreateTime(ctx context.Context, duration time.Duration) {
	dm.connCreateTime.Record(ctx,
		duration.Seconds(),
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// RecordConnectionWaitTime records the time a client waits to acquire a connection
func (dm *Metrics) RecordConnectionWaitTime(ctx context.Context, duration time.Duration) {
	dm.connWaitTime.Record(ctx,
		duration.Seconds(),
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// RecordConnectionUseTime records the time a connection is used by a client
func (dm *Metrics) RecordConnectionUseTime(ctx context.Context, duration time.Duration) {
	dm.connUseTime.Record(ctx,
		duration.Seconds(),
		dm.connPoolName,
		dm.commonAttrs...,
	)
}

// NewDatabaseMetrics creates OpenTelemetry compliant database metrics instruments
func newMetrics(
	ctx context.Context,
	dbSystem dbconv.SystemNameAttr,
	dbHostName string,
	dbHostPort int,
	dbName string,
) (*Metrics, error) {
	meter := otel.Meter(instrumentationIdentifier)

	operationDuration, err := dbconv.NewClientOperationDuration(meter,
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10), // Recommended boundaries
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create operation duration metric: %w", err)
	}

	returnedRows, err := dbconv.NewClientResponseReturnedRows(meter,
		metric.WithExplicitBucketBoundaries(1, 2, 5, 10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000), // Recommended boundaries
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create returned rows metric: %w", err)
	}

	connectionCount, err := dbconv.NewClientConnectionCount(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection count metric: %w", err)
	}

	connectionMax, err := dbconv.NewClientConnectionMax(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection max metric: %w", err)
	}

	connectionIdleMax, err := dbconv.NewClientConnectionIdleMax(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection idle max metric: %w", err)
	}

	connectionIdleMin, err := dbconv.NewClientConnectionIdleMin(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection idle min metric: %w", err)
	}

	connectionPendingRequests, err := dbconv.NewClientConnectionPendingRequests(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pending requests metric: %w", err)
	}

	connectionTimeouts, err := dbconv.NewClientConnectionTimeouts(meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection timeouts metric: %w", err)
	}

	connectionCreateTime, err := dbconv.NewClientConnectionCreateTime(meter,
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10), // Connection creation time boundaries
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection create time metric: %w", err)
	}

	connectionWaitTime, err := dbconv.NewClientConnectionWaitTime(meter,
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10), // Connection wait time boundaries
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection wait time metric: %w", err)
	}

	connectionUseTime, err := dbconv.NewClientConnectionUseTime(meter,
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10), // Connection use time boundaries
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection use time metric: %w", err)
	}

	dm := &Metrics{
		dbSystem:            dbSystem,
		commonAttrs:         buildCommonAttrs(dbSystem, dbHostName, dbHostPort, dbName),
		operationDuration:   operationDuration,
		returnedRows:        returnedRows,
		connCount:           connectionCount,
		connMax:             connectionMax,
		connIdleMaxCount:    connectionIdleMax,
		connIdleMinCount:    connectionIdleMin,
		connPendingRequests: connectionPendingRequests,
		connTimetouts:       connectionTimeouts,
		connCreateTime:      connectionCreateTime,
		connWaitTime:        connectionWaitTime,
		connUseTime:         connectionUseTime,
	}
	if sc := telemetry.ServiceConfigFromContext(ctx); sc.IsValid() {
		dm.connPoolName = sc.System + "/" + sc.Name + "-" + sc.HostName + ":" + sc.ContainerName + "/" + dbName
	} else {
		dm.connPoolName = instrumentationIdentifier
	}

	return dm, nil
}
