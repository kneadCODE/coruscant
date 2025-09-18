package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToOTELAttributes(t *testing.T) {
	tests := []struct {
		name           string
		input          []any
		expectedLength int
		shouldBeNil    bool
	}{
		{
			name: "comprehensive type coverage",
			input: []any{
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
			},
			expectedLength: 15,
			shouldBeNil:    false,
		},
		{
			name:           "empty input",
			input:          []any{},
			expectedLength: 0,
			shouldBeNil:    true,
		},
		{
			name:           "odd number of arguments",
			input:          []any{"key1", "value1", "key2"},
			expectedLength: 0,
			shouldBeNil:    true,
		},
		{
			name:           "invalid key type",
			input:          []any{123, "value", "valid_key", "value"},
			expectedLength: 1,
			shouldBeNil:    false,
		},
		{
			name:           "nil values",
			input:          []any{"nil_key", nil},
			expectedLength: 1,
			shouldBeNil:    false,
		},
		{
			name:           "mixed valid and invalid keys",
			input:          []any{"valid", "value1", 123, "invalid", "another_valid", "value2"},
			expectedLength: 2,
			shouldBeNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToOTELAttributes(tt.input)

			if tt.shouldBeNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedLength)
			}
		})
	}
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
