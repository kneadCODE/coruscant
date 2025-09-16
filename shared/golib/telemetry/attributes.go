package telemetry

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
)

// convertToOTELAttributes converts key-value pairs to OTEL attributes
// This function consolidates the repeated attribute conversion logic used across metric functions.
func convertToOTELAttributes(attrs []any) []attribute.KeyValue {
	if len(attrs) == 0 || len(attrs)%2 != 0 {
		return nil
	}

	otelAttrs := make([]attribute.KeyValue, 0, len(attrs)/2)
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
		case float32:
			otelAttrs = append(otelAttrs, attribute.Float64(key, float64(v)))
		case float64:
			otelAttrs = append(otelAttrs, attribute.Float64(key, v))
		case bool:
			otelAttrs = append(otelAttrs, attribute.Bool(key, v))
		default:
			otelAttrs = append(otelAttrs, attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
	return otelAttrs
}

func replaceAttrKeyDotWithUnderscore(k attribute.Key) string {
	return strings.ReplaceAll(string(k), ".", "_")
}
