package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemID_NewItemID(t *testing.T) {
	id, err := NewItemID()

	require.NoError(t, err)
	assert.True(t, id.IsValid())
	assert.NotEmpty(t, id.String())
}

func TestItemID_NewItemIDFromString(t *testing.T) {
	validID, err := NewItemID()
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
			id, err := NewItemIDFromString(tt.idStr)

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

func TestItemID_Equals(t *testing.T) {
	id1, err := NewItemID()
	require.NoError(t, err)

	id2, err := NewItemID()
	require.NoError(t, err)

	id1Copy, err := NewItemIDFromString(id1.String())
	require.NoError(t, err)

	tests := []struct {
		name     string
		id1      ItemID
		id2      ItemID
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

func TestItemType_NewItemType(t *testing.T) {
	tests := []struct {
		name     string
		itemType string
		want     ItemType
		wantErr  bool
	}{
		{
			name:     "valid income type",
			itemType: "INCOME",
			want:     ItemTypeIncome,
			wantErr:  false,
		},
		{
			name:     "valid expense type",
			itemType: "EXPENSE",
			want:     ItemTypeExpense,
			wantErr:  false,
		},
		{
			name:     "valid transfer type",
			itemType: "TRANSFER",
			want:     ItemTypeTransfer,
			wantErr:  false,
		},
		{
			name:     "invalid type",
			itemType: "INVALID",
			want:     ItemTypeExpense,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewItemType(tt.itemType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestItemType_String(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		want     string
	}{
		{
			name:     "income type string",
			itemType: ItemTypeIncome,
			want:     "INCOME",
		},
		{
			name:     "expense type string",
			itemType: ItemTypeExpense,
			want:     "EXPENSE",
		},
		{
			name:     "transfer type string",
			itemType: ItemTypeTransfer,
			want:     "TRANSFER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.itemType.String())
		})
	}
}

func TestItemType_TypeChecks(t *testing.T) {
	tests := []struct {
		name       string
		itemType   ItemType
		isIncome   bool
		isExpense  bool
		isTransfer bool
	}{
		{
			name:       "income type",
			itemType:   ItemTypeIncome,
			isIncome:   true,
			isExpense:  false,
			isTransfer: false,
		},
		{
			name:       "expense type",
			itemType:   ItemTypeExpense,
			isIncome:   false,
			isExpense:  true,
			isTransfer: false,
		},
		{
			name:       "transfer type",
			itemType:   ItemTypeTransfer,
			isIncome:   false,
			isExpense:  false,
			isTransfer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isIncome, tt.itemType.IsIncome())
			assert.Equal(t, tt.isExpense, tt.itemType.IsExpense())
			assert.Equal(t, tt.isTransfer, tt.itemType.IsTransfer())
		})
	}
}

func TestItemCategory_NewItemCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     ItemCategory
		wantErr  bool
	}{
		{
			name:     "valid salary category",
			category: "SALARY",
			want:     ItemCategorySalary,
			wantErr:  false,
		},
		{
			name:     "valid housing category",
			category: "HOUSING",
			want:     ItemCategoryHousing,
			wantErr:  false,
		},
		{
			name:     "valid savings transfer category",
			category: "SAVINGS_TRANSFER",
			want:     ItemCategorySavingsTransfer,
			wantErr:  false,
		},
		{
			name:     "invalid category",
			category: "INVALID",
			want:     ItemCategoryOtherExpense,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewItemCategory(tt.category)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestItemCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category ItemCategory
		want     string
	}{
		{
			name:     "salary category string",
			category: ItemCategorySalary,
			want:     "SALARY",
		},
		{
			name:     "housing category string",
			category: ItemCategoryHousing,
			want:     "HOUSING",
		},
		{
			name:     "savings transfer category string",
			category: ItemCategorySavingsTransfer,
			want:     "SAVINGS_TRANSFER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.category.String())
		})
	}
}

func TestItemCategory_GetItemType(t *testing.T) {
	tests := []struct {
		name     string
		category ItemCategory
		want     ItemType
	}{
		{
			name:     "salary is income",
			category: ItemCategorySalary,
			want:     ItemTypeIncome,
		},
		{
			name:     "housing is expense",
			category: ItemCategoryHousing,
			want:     ItemTypeExpense,
		},
		{
			name:     "savings transfer is transfer",
			category: ItemCategorySavingsTransfer,
			want:     ItemTypeTransfer,
		},
		{
			name:     "freelance is income",
			category: ItemCategoryFreelance,
			want:     ItemTypeIncome,
		},
		{
			name:     "food is expense",
			category: ItemCategoryFood,
			want:     ItemTypeExpense,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.category.GetItemType())
		})
	}
}

func TestGetIncomeCategories(t *testing.T) {
	categories := GetIncomeCategories()

	assert.Len(t, categories, 8)
	assert.Contains(t, categories, ItemCategorySalary)
	assert.Contains(t, categories, ItemCategoryFreelance)
	assert.Contains(t, categories, ItemCategoryInvestmentReturns)
	assert.Contains(t, categories, ItemCategoryRentalIncome)
	assert.Contains(t, categories, ItemCategorySideBusiness)
	assert.Contains(t, categories, ItemCategoryBonus)
	assert.Contains(t, categories, ItemCategoryGifts)
	assert.Contains(t, categories, ItemCategoryOtherIncome)
}

func TestGetExpenseCategories(t *testing.T) {
	categories := GetExpenseCategories()

	assert.Len(t, categories, 13)
	assert.Contains(t, categories, ItemCategoryHousing)
	assert.Contains(t, categories, ItemCategoryFood)
	assert.Contains(t, categories, ItemCategoryTransportation)
	assert.Contains(t, categories, ItemCategoryUtilities)
	assert.Contains(t, categories, ItemCategoryInsurance)
	assert.Contains(t, categories, ItemCategoryHealthcare)
	assert.Contains(t, categories, ItemCategoryEntertainment)
	assert.Contains(t, categories, ItemCategoryEducation)
	assert.Contains(t, categories, ItemCategoryShopping)
	assert.Contains(t, categories, ItemCategoryPersonalCare)
	assert.Contains(t, categories, ItemCategorySubscriptions)
	assert.Contains(t, categories, ItemCategoryTaxes)
	assert.Contains(t, categories, ItemCategoryOtherExpense)
}

func TestGetTransferCategories(t *testing.T) {
	categories := GetTransferCategories()

	assert.Len(t, categories, 6)
	assert.Contains(t, categories, ItemCategorySavingsTransfer)
	assert.Contains(t, categories, ItemCategoryDebtPayment)
	assert.Contains(t, categories, ItemCategoryInvestmentContribution)
	assert.Contains(t, categories, ItemCategoryEmergencyFund)
	assert.Contains(t, categories, ItemCategoryRetirementContribution)
	assert.Contains(t, categories, ItemCategoryOtherTransfer)
}