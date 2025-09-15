package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToOTELAttributes(t *testing.T) {
	t.Run("comprehensive type coverage", func(t *testing.T) {
		attrs := []any{
			"string_key", "string_value",
			"int_key", 42,
			"int8_key", int8(8),
			"int16_key", int16(16),
			"int32_key", int32(32),
			"int64_key", int64(123),
			"uint_key", uint(42),
			"uint8_key", uint8(8),
			"uint16_key", uint16(16),
			"uint32_key", uint32(32),
			"uint64_key", uint64(64),
			"float32_key", float32(3.14),
			"float64_key", 2.71,
			"bool_key", true,
			"custom_type",
			struct{ Name string }{"test"},
		}

		result := convertToOTELAttributes(attrs)
		assert.Len(t, result, 15) // All 15 key-value pairs should be converted

		// Verify that attributes were created (detailed type checking would be complex)
		assert.NotEmpty(t, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result := convertToOTELAttributes([]any{})
		assert.Nil(t, result)
	})

	t.Run("odd number of arguments", func(t *testing.T) {
		attrs := []any{"key1", "value1", "key2"}
		result := convertToOTELAttributes(attrs)
		assert.Nil(t, result)
	})

	t.Run("invalid key type", func(t *testing.T) {
		attrs := []any{123, "value", "valid_key", "value"}
		result := convertToOTELAttributes(attrs)
		assert.Len(t, result, 1) // Only valid pair should be included
	})

	t.Run("nil values", func(t *testing.T) {
		attrs := []any{"nil_key", nil}
		result := convertToOTELAttributes(attrs)
		assert.Len(t, result, 1) // nil should be converted to string
	})
}

func BenchmarkConvertToOTELAttributes(b *testing.B) {
	attrs := []any{
		"method", "GET",
		"status", 200,
		"path", "/api/users",
		"cached", true,
		"duration", 0.125,
		"retries", int64(3),
		"score", float32(98.5),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = convertToOTELAttributes(attrs)
	}
}
