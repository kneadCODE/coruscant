package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterpartyID_NewCounterpartyID(t *testing.T) {
	id, err := NewCounterpartyID()

	require.NoError(t, err)
	assert.True(t, id.IsValid())
	assert.NotEmpty(t, id.String())
}

func TestCounterpartyID_NewCounterpartyIDFromString(t *testing.T) {
	validID, err := NewCounterpartyID()
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
			id, err := NewCounterpartyIDFromString(tt.idStr)

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

func TestCounterpartyID_Equals(t *testing.T) {
	id1, err := NewCounterpartyID()
	require.NoError(t, err)

	id2, err := NewCounterpartyID()
	require.NoError(t, err)

	id1Copy, err := NewCounterpartyIDFromString(id1.String())
	require.NoError(t, err)

	tests := []struct {
		name     string
		id1      CounterpartyID
		id2      CounterpartyID
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

func TestCounterpartyType_NewCounterpartyType(t *testing.T) {
	tests := []struct {
		name            string
		counterpartyType string
		want            CounterpartyType
		wantErr         bool
	}{
		{
			name:            "valid individual type",
			counterpartyType: "INDIVIDUAL",
			want:            CounterpartyTypeIndividual,
			wantErr:         false,
		},
		{
			name:            "valid business type",
			counterpartyType: "BUSINESS",
			want:            CounterpartyTypeBusiness,
			wantErr:         false,
		},
		{
			name:            "valid bank type",
			counterpartyType: "BANK",
			want:            CounterpartyTypeBank,
			wantErr:         false,
		},
		{
			name:            "invalid type",
			counterpartyType: "INVALID",
			want:            "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCounterpartyType(tt.counterpartyType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCounterpartyType_String(t *testing.T) {
	tests := []struct {
		name string
		cType CounterpartyType
		want string
	}{
		{
			name:  "individual type string",
			cType: CounterpartyTypeIndividual,
			want:  "INDIVIDUAL",
		},
		{
			name:  "business type string",
			cType: CounterpartyTypeBusiness,
			want:  "BUSINESS",
		},
		{
			name:  "online service type string",
			cType: CounterpartyTypeOnlineService,
			want:  "ONLINE_SERVICE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cType.String())
		})
	}
}

func TestAllCounterpartyTypes(t *testing.T) {
	types := AllCounterpartyTypes()

	assert.Len(t, types, 11)
	assert.Contains(t, types, CounterpartyTypeIndividual)
	assert.Contains(t, types, CounterpartyTypeBusiness)
	assert.Contains(t, types, CounterpartyTypeOrganization)
	assert.Contains(t, types, CounterpartyTypeGovernment)
	assert.Contains(t, types, CounterpartyTypeUtility)
	assert.Contains(t, types, CounterpartyTypeOnlineService)
	assert.Contains(t, types, CounterpartyTypeRetailer)
	assert.Contains(t, types, CounterpartyTypeInvestment)
	assert.Contains(t, types, CounterpartyTypeInsurance)
	assert.Contains(t, types, CounterpartyTypeEmployer)
	assert.Contains(t, types, CounterpartyTypeBank)
}
