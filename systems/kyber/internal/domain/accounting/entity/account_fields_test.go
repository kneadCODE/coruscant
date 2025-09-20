package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountID_NewAccountID(t *testing.T) {
	id, err := NewAccountID()

	require.NoError(t, err)
	assert.True(t, id.IsValid())
	assert.NotEmpty(t, id.String())
}

func TestAccountID_NewAccountIDFromString(t *testing.T) {
	validID, err := NewAccountID()
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
			id, err := NewAccountIDFromString(tt.idStr)

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

func TestAccountID_Equals(t *testing.T) {
	id1, err := NewAccountID()
	require.NoError(t, err)

	id2, err := NewAccountID()
	require.NoError(t, err)

	id1Copy, err := NewAccountIDFromString(id1.String())
	require.NoError(t, err)

	tests := []struct {
		name     string
		id1      AccountID
		id2      AccountID
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

func TestAccountCategory_NewAccountCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     AccountCategory
		wantErr  bool
	}{
		{
			name:     "valid asset category",
			category: "ASSET",
			want:     AccountCategoryAsset,
			wantErr:  false,
		},
		{
			name:     "valid liability category",
			category: "LIABILITY",
			want:     AccountCategoryLiability,
			wantErr:  false,
		},
		{
			name:     "valid equity category",
			category: "EQUITY",
			want:     AccountCategoryEquity,
			wantErr:  false,
		},
		{
			name:     "invalid category",
			category: "INVALID",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountCategory(tt.category)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccountCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category AccountCategory
		want     string
	}{
		{
			name:     "asset category string",
			category: AccountCategoryAsset,
			want:     "ASSET",
		},
		{
			name:     "liability category string",
			category: AccountCategoryLiability,
			want:     "LIABILITY",
		},
		{
			name:     "equity category string",
			category: AccountCategoryEquity,
			want:     "EQUITY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.category.String())
		})
	}
}

func TestAccountCategory_TypeChecks(t *testing.T) {
	tests := []struct {
		name        string
		category    AccountCategory
		isAsset     bool
		isLiability bool
		isEquity    bool
	}{
		{
			name:        "asset category",
			category:    AccountCategoryAsset,
			isAsset:     true,
			isLiability: false,
			isEquity:    false,
		},
		{
			name:        "liability category",
			category:    AccountCategoryLiability,
			isAsset:     false,
			isLiability: true,
			isEquity:    false,
		},
		{
			name:        "equity category",
			category:    AccountCategoryEquity,
			isAsset:     false,
			isLiability: false,
			isEquity:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isAsset, tt.category.IsAsset())
			assert.Equal(t, tt.isLiability, tt.category.IsLiability())
			assert.Equal(t, tt.isEquity, tt.category.IsEquity())
		})
	}
}

func TestAllAccountCategories(t *testing.T) {
	categories := AllAccountCategories()

	assert.Len(t, categories, 3)
	assert.Contains(t, categories, AccountCategoryAsset)
	assert.Contains(t, categories, AccountCategoryLiability)
	assert.Contains(t, categories, AccountCategoryEquity)
}

func TestAccountType_NewAccountType(t *testing.T) {
	tests := []struct {
		name        string
		accountType string
		want        AccountType
		wantErr     bool
	}{
		{
			name:        "valid checking type",
			accountType: "CHECKING",
			want:        AccountTypeChecking,
			wantErr:     false,
		},
		{
			name:        "valid credit card type",
			accountType: "CREDIT_CARD",
			want:        AccountTypeCreditCard,
			wantErr:     false,
		},
		{
			name:        "invalid type",
			accountType: "INVALID",
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountType(tt.accountType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
func TestAccountType_String(t *testing.T) {
	tests := []struct {
		name        string
		accountType AccountType
		expected    string
	}{
		{
			name:        "checking account string",
			accountType: AccountTypeChecking,
			expected:    "CHECKING",
		},
		{
			name:        "savings account string",
			accountType: AccountTypeSavings,
			expected:    "SAVINGS",
		},
		{
			name:        "credit card string",
			accountType: AccountTypeCreditCard,
			expected:    "CREDIT_CARD",
		},
		{
			name:        "loan string",
			accountType: AccountTypeLoan,
			expected:    "LOAN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.accountType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAccountStatus_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "active status is active",
			status:   AccountStatusActive,
			expected: true,
		},
		{
			name:     "archived status is not active",
			status:   AccountStatusArchived,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsActive()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAccountStatus_IsArchived(t *testing.T) {
	tests := []struct {
		name     string
		status   AccountStatus
		expected bool
	}{
		{
			name:     "archived status is archived",
			status:   AccountStatusArchived,
			expected: true,
		},
		{
			name:     "active status is not archived",
			status:   AccountStatusActive,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsArchived()
			assert.Equal(t, tt.expected, result)
		})
	}
}