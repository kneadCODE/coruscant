package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
)

func TestNewCounterparty(t *testing.T) {
	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	tests := []struct {
		name             string
		ledgerID         entity.LedgerID
		counterpartyName string
		counterpartyType CounterpartyType
		description      string
		wantErr          bool
		errContains      string
	}{
		{
			name:             "valid counterparty creation",
			ledgerID:         ledgerID,
			counterpartyName: "ABC Company",
			counterpartyType: CounterpartyTypeBusiness,
			description:      "Local grocery store",
			wantErr:          false,
		},
		{
			name:             "invalid ledger ID",
			ledgerID:         entity.LedgerID{},
			counterpartyName: "Test Counterparty",
			counterpartyType: CounterpartyTypeIndividual,
			description:      "Test description",
			wantErr:          true,
			errContains:      "ledger ID cannot be empty",
		},
		{
			name:             "empty counterparty name",
			ledgerID:         ledgerID,
			counterpartyName: "",
			counterpartyType: CounterpartyTypeBusiness,
			description:      "Test description",
			wantErr:          true,
			errContains:      "counterparty name cannot be empty",
		},
		{
			name:             "empty description is allowed",
			ledgerID:         ledgerID,
			counterpartyName: "Valid Counterparty",
			counterpartyType: CounterpartyTypeGovernment,
			description:      "",
			wantErr:          false,
		},
		{
			name:             "individual counterparty",
			ledgerID:         ledgerID,
			counterpartyName: "John Doe",
			counterpartyType: CounterpartyTypeIndividual,
			description:      "Freelance contractor",
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterparty, err := NewCounterparty(tt.ledgerID, tt.counterpartyName, tt.counterpartyType, tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, counterparty)
			} else {
				require.NoError(t, err)
				require.NotNil(t, counterparty)

				assert.True(t, counterparty.ID.IsValid())
				assert.Equal(t, tt.ledgerID, counterparty.LedgerID)
				assert.Equal(t, tt.counterpartyName, counterparty.Name)
				assert.Equal(t, tt.counterpartyType, counterparty.Type)
				assert.Equal(t, tt.description, counterparty.Description)
				assert.Equal(t, CounterpartyStatusActive, counterparty.Status)
				assert.False(t, counterparty.CreatedAt.IsZero())
				assert.False(t, counterparty.UpdatedAt.IsZero())
			}
		})
	}
}

func TestReconstructCounterparty(t *testing.T) {
	counterpartyID, err := NewCounterpartyID()
	require.NoError(t, err)

	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	counterparty := ReconstructCounterparty(
		counterpartyID,
		ledgerID,
		"Reconstructed Counterparty",
		CounterpartyTypeOrganization,
		"Test description",
		"contact@example.com", // This is unused in current implementation but matches signature
		CounterpartyStatusArchived,
		createdAt,
		updatedAt,
	)

	assert.Equal(t, counterpartyID, counterparty.ID)
	assert.Equal(t, ledgerID, counterparty.LedgerID)
	assert.Equal(t, "Reconstructed Counterparty", counterparty.Name)
	assert.Equal(t, CounterpartyTypeOrganization, counterparty.Type)
	assert.Equal(t, "Test description", counterparty.Description)
	assert.Equal(t, CounterpartyStatusArchived, counterparty.Status)
	assert.Equal(t, createdAt, counterparty.CreatedAt)
	assert.Equal(t, updatedAt, counterparty.UpdatedAt)
}

func TestCounterparty_UpdateInfo(t *testing.T) {
	counterparty := createTestCounterparty(t)
	originalUpdatedAt := counterparty.UpdatedAt

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		newName     string
		description string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid update",
			newName:     "Updated Counterparty Name",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "empty name",
			newName:     "",
			description: "Some description",
			wantErr:     true,
			errContains: "counterparty name cannot be empty",
		},
		{
			name:        "empty description is allowed",
			newName:     "Counterparty Name",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := counterparty.UpdateInfo(tt.newName, tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newName, counterparty.Name)
				assert.Equal(t, tt.description, counterparty.Description)
				assert.True(t, counterparty.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestCounterparty_UpdateType(t *testing.T) {
	counterparty := createTestCounterparty(t)
	originalType := counterparty.Type
	originalUpdatedAt := counterparty.UpdatedAt

	time.Sleep(time.Millisecond)

	newType := CounterpartyTypeOrganization
	assert.NotEqual(t, originalType, newType) // Ensure we're actually changing it

	counterparty.UpdateType(newType)

	assert.Equal(t, newType, counterparty.Type)
	assert.True(t, counterparty.UpdatedAt.After(originalUpdatedAt))
}

func TestCounterparty_Activate(t *testing.T) {
	counterparty := createTestCounterparty(t)

	// First archive it
	counterparty.Ardchive() // Note: using the actual method name from the code (with typo)
	assert.Equal(t, CounterpartyStatusArchived, counterparty.Status)

	originalUpdatedAt := counterparty.UpdatedAt
	time.Sleep(time.Millisecond)

	counterparty.Activate()

	assert.Equal(t, CounterpartyStatusActive, counterparty.Status)
	assert.True(t, counterparty.UpdatedAt.After(originalUpdatedAt))
}

func TestCounterparty_Ardchive(t *testing.T) { // Note: testing the actual method name (with typo)
	counterparty := createTestCounterparty(t)
	assert.Equal(t, CounterpartyStatusActive, counterparty.Status)

	originalUpdatedAt := counterparty.UpdatedAt
	time.Sleep(time.Millisecond)

	counterparty.Ardchive() // Note: using the actual method name from the code (with typo)

	assert.Equal(t, CounterpartyStatusArchived, counterparty.Status)
	assert.True(t, counterparty.UpdatedAt.After(originalUpdatedAt))
}

func TestCounterparty_StatusOperations(t *testing.T) {
	counterparty := createTestCounterparty(t)

	// Test initial state
	assert.Equal(t, CounterpartyStatusActive, counterparty.Status)

	// Test archive -> activate cycle
	counterparty.Ardchive()
	assert.Equal(t, CounterpartyStatusArchived, counterparty.Status)

	counterparty.Activate()
	assert.Equal(t, CounterpartyStatusActive, counterparty.Status)

	// Test multiple state changes
	for i := 0; i < 3; i++ {
		counterparty.Ardchive()
		assert.Equal(t, CounterpartyStatusArchived, counterparty.Status)

		counterparty.Activate()
		assert.Equal(t, CounterpartyStatusActive, counterparty.Status)
	}
}

func TestCounterparty_AllCounterpartyTypes(t *testing.T) {
	// Test that we can create counterparties with all supported types
	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	allTypes := AllCounterpartyTypes()
	assert.Len(t, allTypes, 11) // Should have 11 types based on our implementation

	for _, counterpartyType := range allTypes {
		t.Run(string(counterpartyType), func(t *testing.T) {
			counterparty, err := NewCounterparty(
				ledgerID,
				"Test "+string(counterpartyType),
				counterpartyType,
				"Test description for "+string(counterpartyType),
			)

			require.NoError(t, err)
			assert.Equal(t, counterpartyType, counterparty.Type)
		})
	}
}

func TestCounterparty_TypeSpecificBehavior(t *testing.T) {
	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	tests := []struct {
		name             string
		counterpartyType CounterpartyType
		expectedName     string
		description      string
	}{
		{
			name:             "individual counterparty",
			counterpartyType: CounterpartyTypeIndividual,
			expectedName:     "John Smith",
			description:      "Freelance consultant",
		},
		{
			name:             "business counterparty",
			counterpartyType: CounterpartyTypeBusiness,
			expectedName:     "ABC Corp",
			description:      "Software development company",
		},
		{
			name:             "government counterparty",
			counterpartyType: CounterpartyTypeGovernment,
			expectedName:     "Internal Revenue Service",
			description:      "Tax authority",
		},
		{
			name:             "bank counterparty",
			counterpartyType: CounterpartyTypeBank,
			expectedName:     "First National Bank",
			description:      "Primary banking institution",
		},
		{
			name:             "utility counterparty",
			counterpartyType: CounterpartyTypeUtility,
			expectedName:     "City Electric Company",
			description:      "Electricity provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterparty, err := NewCounterparty(
				ledgerID,
				tt.expectedName,
				tt.counterpartyType,
				tt.description,
			)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, counterparty.Name)
			assert.Equal(t, tt.counterpartyType, counterparty.Type)
			assert.Equal(t, tt.description, counterparty.Description)

			// Test type updates
			if tt.counterpartyType != CounterpartyTypeIndividual {
				counterparty.UpdateType(CounterpartyTypeIndividual)
				assert.Equal(t, CounterpartyTypeIndividual, counterparty.Type)
			}
		})
	}
}

// Helper functions

func createTestCounterparty(t *testing.T) *Counterparty {
	t.Helper()

	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	counterparty, err := NewCounterparty(
		ledgerID,
		"Test Counterparty",
		CounterpartyTypeBusiness,
		"Test description",
	)
	require.NoError(t, err)

	return counterparty
}