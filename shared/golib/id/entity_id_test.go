package id

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEntityID(t *testing.T) {
	t.Run("creates valid EntityID", func(t *testing.T) {
		id, err := NewEntityID()
		require.NoError(t, err)
		assert.True(t, id.IsValid())
		assert.NotEmpty(t, id.String())
	})

	t.Run("creates unique IDs", func(t *testing.T) {
		id1, err := NewEntityID()
		require.NoError(t, err)

		id2, err := NewEntityID()
		require.NoError(t, err)

		assert.False(t, id1.Equals(id2))
		assert.NotEqual(t, id1.String(), id2.String())
	})

	t.Run("creates UUIDv7 format", func(t *testing.T) {
		id, err := NewEntityID()
		require.NoError(t, err)

		// UUIDv7 should be 36 characters with dashes
		idStr := id.String()
		assert.Len(t, idStr, 36)
		assert.Equal(t, 4, strings.Count(idStr, "-"))

		// Should contain timestamp that is recent
		timestamp, err := id.Timestamp()
		require.NoError(t, err)
		assert.True(t, time.Since(timestamp) < time.Minute)
	})
}

func TestNewEntityIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid UUID string",
			input:   "018c2a2a-1234-7abc-8def-0123456789ab",
			wantErr: false,
		},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			errContains: "entity ID cannot be empty",
		},
		{
			name:        "invalid UUID format",
			input:       "not-a-uuid",
			wantErr:     true,
			errContains: "invalid entity ID format",
		},
		{
			name:        "invalid UUID characters",
			input:       "018c2a2a-1234-7abc-8def-0123456789zz",
			wantErr:     true,
			errContains: "invalid entity ID format",
		},
		{
			name:        "wrong length",
			input:       "018c2a2a-1234-7abc-8def-012345",
			wantErr:     true,
			errContains: "invalid entity ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewEntityIDFromString(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, EntityID{}, id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.input, id.String())
				assert.True(t, id.IsValid())
			}
		})
	}
}

func TestEntityID_String(t *testing.T) {
	t.Run("returns correct string representation", func(t *testing.T) {
		testUUID := "018c2a2a-1234-7abc-8def-0123456789ab"
		id, err := NewEntityIDFromString(testUUID)
		require.NoError(t, err)

		assert.Equal(t, testUUID, id.String())
	})

	t.Run("empty EntityID returns empty string", func(t *testing.T) {
		var id EntityID
		assert.Equal(t, "", id.String())
	})
}

func TestEntityID_Equals(t *testing.T) {
	testUUID1 := "018c2a2a-1234-7abc-8def-0123456789ab"
	testUUID2 := "018c2a2a-5678-7def-8abc-0123456789cd"

	tests := []struct {
		name     string
		id1      EntityID
		id2      EntityID
		expected bool
	}{
		{
			name:     "same ID equals itself",
			id1:      mustCreateFromString(t, testUUID1),
			id2:      mustCreateFromString(t, testUUID1),
			expected: true,
		},
		{
			name:     "different IDs not equal",
			id1:      mustCreateFromString(t, testUUID1),
			id2:      mustCreateFromString(t, testUUID2),
			expected: false,
		},
		{
			name:     "empty IDs are equal",
			id1:      EntityID{},
			id2:      EntityID{},
			expected: true,
		},
		{
			name:     "empty vs non-empty not equal",
			id1:      EntityID{},
			id2:      mustCreateFromString(t, testUUID1),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.id1.Equals(tt.id2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEntityID_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		id       EntityID
		expected bool
	}{
		{
			name:     "valid UUID is valid",
			id:       mustCreateFromString(t, "018c2a2a-1234-7abc-8def-0123456789ab"),
			expected: true,
		},
		{
			name:     "empty EntityID is invalid",
			id:       EntityID{},
			expected: false,
		},
		{
			name:     "newly created EntityID is valid",
			id:       mustCreate(t),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.id.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEntityID_Timestamp(t *testing.T) {
	t.Run("valid EntityID returns timestamp", func(t *testing.T) {
		beforeCreation := time.Now()
		id := mustCreate(t)
		afterCreation := time.Now()

		timestamp, err := id.Timestamp()
		require.NoError(t, err)

		// Timestamp should be between before and after creation
		assert.True(t, timestamp.After(beforeCreation.Add(-time.Second))) // Allow 1s margin
		assert.True(t, timestamp.Before(afterCreation.Add(time.Second)))  // Allow 1s margin
	})

	t.Run("empty EntityID returns error", func(t *testing.T) {
		var id EntityID
		_, err := id.Timestamp()
		assert.Error(t, err)
	})

	t.Run("known timestamp extraction", func(t *testing.T) {
		// Use a known UUIDv7 with embedded timestamp for predictable testing
		id, err := NewEntityIDFromString("018c2a2a-1234-7abc-8def-0123456789ab")
		require.NoError(t, err)

		timestamp, err := id.Timestamp()
		require.NoError(t, err)

		// The timestamp embedded in this UUIDv7 corresponds to a specific time
		assert.False(t, timestamp.IsZero())
		assert.True(t, timestamp.Before(time.Now()))
	})
}

func TestEntityID_TypeSafety(t *testing.T) {
	t.Run("EntityID is type-safe", func(t *testing.T) {
		id1 := mustCreate(t)
		id2 := mustCreate(t)

		// This should compile - same types
		assert.False(t, id1.Equals(id2))

		// These would not compile - different types (if we had other ID types)
		// This test verifies type safety conceptually
		assert.IsType(t, EntityID{}, id1)
		assert.IsType(t, EntityID{}, id2)
	})
}

// Helper functions for tests

func mustCreate(t *testing.T) EntityID {
	t.Helper()
	id, err := NewEntityID()
	require.NoError(t, err)
	return id
}

func mustCreateFromString(t *testing.T, s string) EntityID {
	t.Helper()
	id, err := NewEntityIDFromString(s)
	require.NoError(t, err)
	return id
}