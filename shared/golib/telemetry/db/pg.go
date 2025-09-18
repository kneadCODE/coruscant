package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/dbconv"
	"go.opentelemetry.io/otel/trace"
)

// PGXTracker provides automatic observability for PostgreSQL operations using pgx hooks.
// It implements all pgx tracer interfaces to capture metrics and traces without manual instrumentation.
//
// Features:
//   - Automatic OpenTelemetry span creation and management
//   - Database metrics collection (duration, row counts, connection stats)
//   - PostgreSQL semantic convention compliance
//   - Perfect timing alignment between spans and metrics
//   - Zero-configuration observability for all database operations
type PGXTracker struct {
	*Metrics                              // Database metrics collection
	tracer           trace.Tracer         // OpenTelemetry tracer for span creation
	commonTraceAttrs []attribute.KeyValue // Common attributes applied to all spans
}

// NewPGXTracker creates a new PostgreSQL observability tracker.
//
// The tracker automatically instruments all PostgreSQL operations performed through pgx
// by implementing pgx tracer interfaces. It collects both OpenTelemetry metrics and
// distributed traces with full semantic convention compliance.
//
// Parameters:
//   - ctx: Context for initialization (used for telemetry service configuration)
//   - serverAddr: PostgreSQL server address (used in telemetry attributes)
//   - serverPort: PostgreSQL server port (used in telemetry attributes)
//   - database: Database name (used as db.namespace attribute)
//
// Returns a PGXTracker ready to be assigned to pgxpool.Config.ConnConfig.Tracer
func NewPGXTracker(
	ctx context.Context,
	serverAddr string,
	serverPort int,
	database string,
) (*PGXTracker, error) {
	m, err := newMetrics(ctx, dbconv.SystemNamePostgreSQL, serverAddr, serverPort, database)
	if err != nil {
		return nil, err
	}

	return &PGXTracker{
		Metrics: m,
		tracer:  otel.Tracer(instrumentationIdentifier),
		commonTraceAttrs: buildCommonAttrs(
			dbconv.SystemNamePostgreSQL,
			serverAddr,
			serverPort,
			database,
		),
	}, nil
}

// TraceQueryStart implements pgx.QueryTracer interface
func (tr *PGXTracker) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return tr.startMeasuring(ctx, "", data.SQL, true)
}

// TraceQueryEnd implements pgx.QueryTracer interface
func (tr *PGXTracker) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	meta := tr.finishMeasuring(ctx, data.CommandTag.RowsAffected(), data.Err)
	tr.operationDuration.Record(ctx, meta.duration.Seconds(), tr.dbSystem, meta.attrs...)

	// Only record returned rows for operations that actually return data
	if meta.rowsAffected != nil {
		tr.returnedRows.Record(ctx, *meta.rowsAffected, tr.dbSystem, meta.attrs...)
	}
}

// TraceBatchStart implements pgx.BatchTracer interface
func (tr *PGXTracker) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	// Store batch size for metrics
	ctx = context.WithValue(ctx, batchSizeKey{}, data.Batch.Len())

	return tr.startMeasuring(ctx, "BATCH", "", true)
}

// TraceBatchQuery implements pgx.BatchTracer interface
// Intentionally simple - individual queries are tracked as part of the overall batch
func (tr *PGXTracker) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	meta := buildQueryMetadata("", data.SQL)
	meta.attrs = append(meta.attrs, tr.commonTraceAttrs...)

	// Don't record duration here - batch queries are part of the overall batch timing
	// Only record returned rows for operations that actually return data
	if meta.rowsAffected != nil {
		tr.returnedRows.Record(ctx, *meta.rowsAffected, tr.dbSystem, meta.attrs...)
	}

	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.AddEvent("batch.query.executed", trace.WithAttributes(meta.attrs...))
	}
}

// TraceBatchEnd implements pgx.BatchTracer interface
func (tr *PGXTracker) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	// Add batch size attribute to metadata before finishing measuring
	if batchSize, ok := ctx.Value(batchSizeKey{}).(int); ok {
		if meta, metaOk := ctx.Value(pgQueryMetadataKey{}).(pgQueryMetadata); metaOk {
			meta.attrs = append(meta.attrs, semconv.DBOperationBatchSize(batchSize))
			ctx = context.WithValue(ctx, pgQueryMetadataKey{}, meta)
		}
	}

	meta := tr.finishMeasuring(ctx, 0, data.Err)
	tr.operationDuration.Record(ctx, meta.duration.Seconds(), tr.dbSystem, meta.attrs...)
}

// TraceConnectStart implements pgx.ConnectTracer interface
func (tr *PGXTracker) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	return tr.startMeasuring(ctx, "CONNECT", "", false)
}

// TraceConnectEnd implements pgx.ConnectTracer interface
func (tr *PGXTracker) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	startTime, ok := ctx.Value(connectStartTimeKey{}).(time.Time)
	if !ok {
		return // No start time recorded
	}

	duration := time.Since(startTime)

	// Record connection metrics without creating spans
	if data.Err == nil {
		tr.RecordConnectionCreateTime(ctx, duration)
	} else if isConnectionTimeoutError(data.Err) {
		tr.RecordConnectionTimeouts(ctx, 1)
	}
}

// TracePrepareStart implements pgx.PrepareTracer interface
func (tr *PGXTracker) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	return tr.startMeasuring(ctx, "PREPARE", data.SQL, true)
}

// TracePrepareEnd implements pgx.PrepareTracer interface
func (tr *PGXTracker) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
	_ = tr.finishMeasuring(ctx, 0, data.Err)
}

func (tr *PGXTracker) finishMeasuring(ctx context.Context, rowsAffected int64, err error) pgQueryMetadata {
	meta, ok := ctx.Value(pgQueryMetadataKey{}).(pgQueryMetadata)
	if !ok {
		return pgQueryMetadata{}
	}

	if err == nil {
		meta.attrs = append(meta.attrs, semconv.DBResponseStatusCode(pgerrcode.SuccessfulCompletion))
		_, ok := opsNotReturningRows[meta.opName]
		if !ok {
			meta.rowsAffected = &rowsAffected
			meta.attrs = append(meta.attrs, semconv.DBResponseReturnedRows(int(rowsAffected))) // intentionally forcing the conversion
		}
	} else {
		errType, respCode := extractErrTypeAndResultCode(err)
		meta.attrs = append(meta.attrs,
			semconv.ErrorTypeKey.String(errType),
			semconv.DBResponseStatusCode(respCode),
		)
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		// No active span - calculate duration now for metrics
		meta.duration = time.Since(meta.start)
		return meta
	}

	span.SetAttributes(meta.attrs...)
	if err == nil {
		span.SetStatus(codes.Ok, "")
	} else {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	// Calculate duration right before ending span for maximum alignment
	meta.duration = time.Since(meta.start)
	span.End()

	return meta
}

// opName is optional
func (tr *PGXTracker) startMeasuring(ctx context.Context, opName, query string, withSpan bool) context.Context {
	meta := buildQueryMetadata(opName, query)
	meta.attrs = append(meta.attrs, tr.commonTraceAttrs...)

	ctx = context.WithValue(ctx, pgQueryMetadataKey{}, meta)

	if withSpan {
		ctx, _ = tr.tracer.Start(ctx, meta.spanName, trace.WithAttributes(meta.attrs...))
	}
	return ctx
}

// isConnectionTimeoutError checks if an error is related to connection timeout
func isConnectionTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "connection refused")
}

// Context keys for storing data
type (
	pgQueryMetadataKey  struct{}
	batchSizeKey        struct{}
	connectStartTimeKey struct{}
)

type pgQueryMetadata struct {
	attrs        []attribute.KeyValue
	start        time.Time
	opName       string
	spanName     string
	duration     time.Duration
	rowsAffected *int64
}

// Returns empty string if not easily determined
func extractSQLOperationName(query string) string {
	if query == "" {
		return ""
	}

	// Skip leading whitespace
	i := 0
	for i < len(query) && (query[i] == ' ' || query[i] == '\t' || query[i] == '\n') {
		i++
	}

	// Find end of first word
	start := i
	for i < len(query) && query[i] != ' ' && query[i] != '\t' && query[i] != '\n' && query[i] != '(' {
		i++
	}

	if i > start {
		return query[start:i]
	}
	return ""
}

// Returns empty string if not easily determined
func extractSQLCollectionName(query string) string {
	if len(query) < 10 {
		return ""
	}

	// Simple regex-like matching for common patterns
	queryUpper := strings.ToUpper(query)

	// FROM table
	if idx := strings.Index(queryUpper, " FROM "); idx != -1 {
		return extractSQLSimpleTableName(query, idx+6)
	}

	// INTO table
	if idx := strings.Index(queryUpper, " INTO "); idx != -1 {
		return extractSQLSimpleTableName(query, idx+6)
	}

	// UPDATE table
	if strings.HasPrefix(queryUpper, "UPDATE ") {
		return extractSQLSimpleTableName(query, 7)
	}

	return ""
}

// Returns empty string if not easily determined
func extractSQLSimpleTableName(query string, start int) string {
	// Skip whitespace
	for start < len(query) && (query[start] == ' ' || query[start] == '\t') {
		start++
	}

	// Find end of table name
	end := start
	for end < len(query) && query[end] != ' ' && query[end] != '\t' && query[end] != '(' && query[end] != ',' {
		end++
	}

	if end > start && end-start < 50 { // Reasonable table name length
		return query[start:end]
	}
	return ""
}

// returns error type and result code
func extractErrTypeAndResultCode(err error) (string, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code, pgErr.Code
	}

	return string(dbconv.ErrorTypeOther), "ERROR"
}

// opName is optional
func buildQueryMetadata(opName, query string) pgQueryMetadata {
	// attrs := tr.commonAttrs

	var attrs []attribute.KeyValue

	opName = strings.ToUpper(opName)

	var tableName string
	if query != "" {
		attrs = append(attrs, semconv.DBQueryText(query))
		if opName == "" {
			opName = strings.ToUpper(extractSQLOperationName(query))
		}
		tableName = extractSQLCollectionName(query)
	}
	spanName := opName

	if tableName != "" {
		attrs = append(attrs, semconv.DBCollectionName(tableName))
	}
	if opName != "" {
		attrs = append(attrs, semconv.DBOperationName(opName))
		if tableName != "" {
			spanName = opName + " " + tableName
			attrs = append(attrs, semconv.DBQuerySummary(spanName))
		} else {
			attrs = append(attrs, semconv.DBQuerySummary(opName))
		}
	}

	return pgQueryMetadata{
		attrs:    attrs,
		start:    time.Now(),
		spanName: spanName,
		opName:   opName,
	}
}

// Operations that do not have rows returned
var opsNotReturningRows = map[string]bool{
	"BEGIN":    true,
	"COMMIT":   true,
	"ROLLBACK": true,
	"CONNECT":  true,
	"PREPARE":  true,
}

// Ensure we implement all required tracer interfaces
var (
	_ pgx.QueryTracer   = (*PGXTracker)(nil)
	_ pgx.BatchTracer   = (*PGXTracker)(nil)
	_ pgx.ConnectTracer = (*PGXTracker)(nil)
	_ pgx.PrepareTracer = (*PGXTracker)(nil)
)
