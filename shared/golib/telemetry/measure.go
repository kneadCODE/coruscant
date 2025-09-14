package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Measure starts a new span for measuring operation duration and attributes.
// Returns a context with the span and a finish function to end the span.
// Attributes are provided as key-value pairs where keys must be strings.
// Supported value types: string, int, int64, float64, bool. Other types are converted to strings.
// Usage: ctx, finish := Measure(ctx, "operation-name", "key1", value1, "key2", value2)
func Measure(ctx context.Context, operationName string, attrs ...any) (context.Context, func(error)) {
	// Get tracer from the global provider - this is the OTEL recommended approach
	tracer := otel.Tracer(instrumentationIdentifier)

	// Configure span start options
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindInternal),
	}

	// Add attributes if provided (convert key-value pairs to OTEL attributes)
	if len(attrs) > 0 && len(attrs)%2 == 0 {
		var otelAttrs []attribute.KeyValue
		for i := 0; i < len(attrs); i += 2 {
			key, ok := attrs[i].(string)
			if !ok {
				continue // Skip invalid key (must be string)
			}

			value := attrs[i+1]
			switch v := value.(type) {
			case string:
				otelAttrs = append(otelAttrs, attribute.String(key, v))
			case int:
				otelAttrs = append(otelAttrs, attribute.Int(key, v))
			case int64:
				otelAttrs = append(otelAttrs, attribute.Int64(key, v))
			case float64:
				otelAttrs = append(otelAttrs, attribute.Float64(key, v))
			case bool:
				otelAttrs = append(otelAttrs, attribute.Bool(key, v))
			default:
				// For other types, convert to string
				otelAttrs = append(otelAttrs, attribute.String(key, fmt.Sprintf("%v", v)))
			}
		}
		opts = append(opts, trace.WithAttributes(otelAttrs...))
	}

	// Start the span - OTEL will automatically handle trace context propagation
	ctx, span := tracer.Start(ctx, operationName, opts...)

	return ctx, func(err error) {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}
}
