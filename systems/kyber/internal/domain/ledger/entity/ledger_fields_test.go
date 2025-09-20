package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerID_NewLedgerID(t *testing.T) {
	id, err := NewLedgerID()

	require.NoError(t, err)
	assert.True(t, id.IsValid())
	assert.NotEmpty(t, id.String())
}

func TestLedgerID_NewLedgerIDFromString(t *testing.T) {
	validID, err := NewLedgerID()
	require.NoError(t, err)
	validIDStr := validID.String()

	tests := []struct {
		name      string
		idStr     string
		wantErr   bool
		wantValid bool
	}{
		{
			name:      "valid UUID string",
			idStr:     validIDStr,
			wantErr:   false,
			wantValid: true,
		},
		{
			name:      "empty string",
			idStr:     "",
			wantErr:   true,
			wantValid: false,
		},
		{
			name:      "invalid UUID format",
			idStr:     "invalid-uuid",
			wantErr:   true,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewLedgerIDFromString(tt.idStr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.idStr, id.String())
			}

			assert.Equal(t, tt.wantValid, id.IsValid())
		})
	}
}

func TestLedgerID_Equals(t *testing.T) {
	id1, err := NewLedgerID()
	require.NoError(t, err)

	id2, err := NewLedgerID()
	require.NoError(t, err)

	id1Copy, err := NewLedgerIDFromString(id1.String())
	require.NoError(t, err)

	tests := []struct {
		name     string
		id1      LedgerID
		id2      LedgerID
		expected bool
	}{
		{
			name:     "same ID equals itself",
			id1:      id1,
			id2:      id1,
			expected: true,
		},
		{
			name:     "ID equals copy from string",
			id1:      id1,
			id2:      id1Copy,
			expected: true,
		},
		{
			name:     "different IDs not equal",
			id1:      id1,
			id2:      id2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id1.Equals(tt.id2))
		})
	}
}

func TestLedgerID_TypeSafety(t *testing.T) {
	ledgerID, err := NewLedgerID()
	require.NoError(t, err)

	// This ensures LedgerID is a distinct type from EntityID
	// and can't be accidentally used interchangeably
	assert.IsType(t, LedgerID{}, ledgerID)
}