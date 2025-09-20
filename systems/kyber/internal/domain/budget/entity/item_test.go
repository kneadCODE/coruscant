package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

func TestNewItem(t *testing.T) {
	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	tests := []struct {
		name        string
		ledgerID    entity.LedgerID
		itemName    string
		description string
		itemType    ItemType
		currency    money.Currency
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid item creation",
			ledgerID:    ledgerID,
			itemName:    "Groceries",
			description: "Monthly grocery budget",
			itemType:    ItemTypeExpense,
			currency:    "USD",
			wantErr:     false,
		},
		{
			name:        "invalid ledger ID",
			ledgerID:    entity.LedgerID{},
			itemName:    "Test Item",
			description: "Test description",
			itemType:    ItemTypeIncome,
			currency:    "USD",
			wantErr:     true,
			errContains: "ledger ID cannot be empty",
		},
		{
			name:        "empty item name",
			ledgerID:    ledgerID,
			itemName:    "",
			description: "Test description",
			itemType:    ItemTypeExpense,
			currency:    "USD",
			wantErr:     true,
			errContains: "item name cannot be empty",
		},
		{
			name:        "empty currency",
			ledgerID:    ledgerID,
			itemName:    "Test Item",
			description: "Test description",
			itemType:    ItemTypeIncome,
			currency:    "",
			wantErr:     true,
			errContains: "item currency cannot be empty",
		},
		{
			name:        "valid income item",
			ledgerID:    ledgerID,
			itemName:    "Salary",
			description: "Monthly salary",
			itemType:    ItemTypeIncome,
			currency:    "USD",
			wantErr:     false,
		},
		{
			name:        "valid transfer item",
			ledgerID:    ledgerID,
			itemName:    "Savings Transfer",
			description: "Transfer to savings",
			itemType:    ItemTypeTransfer,
			currency:    "USD",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewItem(tt.ledgerID, tt.itemName, tt.description, tt.itemType, tt.currency)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, item)
			} else {
				require.NoError(t, err)
				require.NotNil(t, item)

				assert.True(t, item.ID.IsValid())
				assert.Equal(t, tt.ledgerID, item.LedgerID)
				assert.Equal(t, tt.itemName, item.Name)
				assert.Equal(t, tt.description, item.Description)
				assert.Equal(t, tt.itemType, item.Type)
				assert.Equal(t, tt.currency, item.Currency)
				assert.True(t, item.IsActive)
				assert.NotNil(t, item.MonthlyBudgets)
				assert.Empty(t, item.MonthlyBudgets)
				assert.False(t, item.CreatedAt.IsZero())
				assert.False(t, item.UpdatedAt.IsZero())
			}
		})
	}
}

func TestReconstructItem(t *testing.T) {
	itemID, err := NewItemID()
	require.NoError(t, err)

	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	monthlyBudgets := map[string]*BudgetTracking{
		"2024-01": {
			Year:           2024,
			Month:          1,
			TargetAmount:   mustMoney(t, "100.00", "USD"),
			BudgetedAmount: mustMoney(t, "100.00", "USD"),
			ActualAmount:   mustMoney(t, "75.50", "USD"),
			UpdatedAt:      updatedAt,
		},
	}

	tests := []struct {
		name           string
		monthlyBudgets map[string]*BudgetTracking
	}{
		{
			name:           "with monthly budgets",
			monthlyBudgets: monthlyBudgets,
		},
		{
			name:           "with nil monthly budgets",
			monthlyBudgets: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := ReconstructItem(
				itemID,
				ledgerID,
				"Reconstructed Item",
				"Test description",
				ItemTypeExpense,
				"USD",
				tt.monthlyBudgets,
				false,
				createdAt,
				updatedAt,
			)

			assert.Equal(t, itemID, item.ID)
			assert.Equal(t, ledgerID, item.LedgerID)
			assert.Equal(t, "Reconstructed Item", item.Name)
			assert.Equal(t, "Test description", item.Description)
			assert.Equal(t, ItemTypeExpense, item.Type)
			assert.Equal(t, money.Currency("USD"), item.Currency)
			assert.False(t, item.IsActive)
			assert.Equal(t, createdAt, item.CreatedAt)
			assert.Equal(t, updatedAt, item.UpdatedAt)
			assert.NotNil(t, item.MonthlyBudgets)

			if tt.monthlyBudgets == nil {
				assert.Empty(t, item.MonthlyBudgets)
			} else {
				assert.Equal(t, tt.monthlyBudgets, item.MonthlyBudgets)
			}
		})
	}
}

func TestItem_UpdateInfo(t *testing.T) {
	item := createTestItem(t)
	originalUpdatedAt := item.UpdatedAt

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
			newName:     "Updated Item Name",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "empty name",
			newName:     "",
			description: "Some description",
			wantErr:     true,
			errContains: "item name cannot be empty",
		},
		{
			name:        "empty description is allowed",
			newName:     "Item Name",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := item.UpdateInfo(tt.newName, tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newName, item.Name)
				assert.Equal(t, tt.description, item.Description)
				assert.True(t, item.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestItem_SetMonthlyTarget(t *testing.T) {
	item := createTestItem(t)

	tests := []struct {
		name         string
		year         int
		month        int
		targetAmount money.Money
		wantErr      bool
		errContains  string
	}{
		{
			name:         "set new monthly target",
			year:         2024,
			month:        3,
			targetAmount: mustMoney(t, "500.00", "USD"),
			wantErr:      false,
		},
		{
			name:         "currency mismatch",
			year:         2024,
			month:        4,
			targetAmount: mustMoney(t, "300.00", "EUR"),
			wantErr:      true,
			errContains:  "currency mismatch",
		},
		{
			name:         "update existing monthly target",
			year:         2024,
			month:        3, // Same as first test
			targetAmount: mustMoney(t, "600.00", "USD"),
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := item.SetMonthlyTarget(tt.year, tt.month, tt.targetAmount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)

				monthKey := "2024-03"
				if tt.month == 4 {
					monthKey = "2024-04"
				}

				budget := item.MonthlyBudgets[monthKey]
				require.NotNil(t, budget)
				assert.Equal(t, tt.targetAmount, budget.TargetAmount)
				assert.Equal(t, tt.targetAmount, budget.BudgetedAmount)
			}
		})
	}
}

func TestItem_UpdateMonthlyBudget(t *testing.T) {
	item := createTestItem(t)

	// First set a target
	targetAmount := mustMoney(t, "500.00", "USD")
	err := item.SetMonthlyTarget(2024, 3, targetAmount)
	require.NoError(t, err)

	tests := []struct {
		name           string
		year           int
		month          int
		budgetedAmount money.Money
		wantErr        bool
		errContains    string
	}{
		{
			name:           "update existing budget",
			year:           2024,
			month:          3,
			budgetedAmount: mustMoney(t, "450.00", "USD"),
			wantErr:        false,
		},
		{
			name:           "update non-existent budget",
			year:           2024,
			month:          5,
			budgetedAmount: mustMoney(t, "300.00", "USD"),
			wantErr:        true,
			errContains:    "no budget tracking found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := item.UpdateMonthlyBudget(tt.year, tt.month, tt.budgetedAmount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)

				monthKey := "2024-03"
				budget := item.MonthlyBudgets[monthKey]
				require.NotNil(t, budget)
				assert.Equal(t, tt.budgetedAmount, budget.BudgetedAmount)
			}
		})
	}
}

func TestItem_AddActualAmount(t *testing.T) {
	item := createTestItem(t)

	tests := []struct {
		name    string
		year    int
		month   int
		amount  money.Money
		wantErr bool
	}{
		{
			name:    "add to non-existent budget (creates new)",
			year:    2024,
			month:   6,
			amount:  mustMoney(t, "75.00", "USD"),
			wantErr: false,
		},
		{
			name:    "add to existing budget",
			year:    2024,
			month:   6,
			amount:  mustMoney(t, "25.00", "USD"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := item.AddActualAmount(tt.year, tt.month, tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				monthKey := "2024-06"
				budget := item.MonthlyBudgets[monthKey]
				require.NotNil(t, budget)
			}
		})
	}

	// Verify total actual amount
	budget := item.MonthlyBudgets["2024-06"]
	expectedTotal := mustMoney(t, "100.00", "USD") // 75 + 25
	assert.Equal(t, expectedTotal, budget.ActualAmount)
}

func TestItem_GetMonthlyBudget(t *testing.T) {
	item := createTestItem(t)

	// Set up a budget for testing
	targetAmount := mustMoney(t, "300.00", "USD")
	err := item.SetMonthlyTarget(2024, 7, targetAmount)
	require.NoError(t, err)

	tests := []struct {
		name        string
		year        int
		month       int
		expectFound bool
	}{
		{
			name:        "get existing budget",
			year:        2024,
			month:       7,
			expectFound: true,
		},
		{
			name:        "get non-existent budget",
			year:        2024,
			month:       8,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			budget := item.GetMonthlyBudget(tt.year, tt.month)

			if tt.expectFound {
				require.NotNil(t, budget)
				assert.Equal(t, tt.year, budget.Year)
				assert.Equal(t, tt.month, budget.Month)
			} else {
				assert.Nil(t, budget)
			}
		})
	}
}

func TestItem_GetCurrentMonthBudget(t *testing.T) {
	item := createTestItem(t)

	now := time.Now()
	targetAmount := mustMoney(t, "400.00", "USD")
	err := item.SetMonthlyTarget(now.Year(), int(now.Month()), targetAmount)
	require.NoError(t, err)

	budget := item.GetCurrentMonthBudget()
	require.NotNil(t, budget)
	assert.Equal(t, now.Year(), budget.Year)
	assert.Equal(t, int(now.Month()), budget.Month)
}

func TestItem_ActivateAndDeactivate(t *testing.T) {
	item := createTestItem(t)
	originalUpdatedAt := item.UpdatedAt

	// Test deactivate
	time.Sleep(time.Millisecond)
	item.Deactivate()
	assert.False(t, item.IsActive)
	assert.True(t, item.UpdatedAt.After(originalUpdatedAt))

	// Test activate
	newUpdatedAt := item.UpdatedAt
	time.Sleep(time.Millisecond)
	item.Activate()
	assert.True(t, item.IsActive)
	assert.True(t, item.UpdatedAt.After(newUpdatedAt))
}

func TestItem_GetTotalBudgetedForYear(t *testing.T) {
	item := createTestItem(t)

	// Set up budgets for different months in 2024
	err := item.SetMonthlyTarget(2024, 1, mustMoney(t, "100.00", "USD"))
	require.NoError(t, err)
	err = item.SetMonthlyTarget(2024, 2, mustMoney(t, "150.00", "USD"))
	require.NoError(t, err)
	err = item.SetMonthlyTarget(2024, 3, mustMoney(t, "200.00", "USD"))
	require.NoError(t, err)

	// Set up a budget for a different year
	err = item.SetMonthlyTarget(2023, 12, mustMoney(t, "50.00", "USD"))
	require.NoError(t, err)

	total, err := item.GetTotalBudgetedForYear(2024)
	require.NoError(t, err)

	expectedTotal := mustMoney(t, "450.00", "USD") // 100 + 150 + 200
	assert.Equal(t, expectedTotal, total)

	// Test different year
	total2023, err := item.GetTotalBudgetedForYear(2023)
	require.NoError(t, err)

	expectedTotal2023 := mustMoney(t, "50.00", "USD")
	assert.Equal(t, expectedTotal2023, total2023)

	// Test year with no budgets
	totalEmpty, err := item.GetTotalBudgetedForYear(2025)
	require.NoError(t, err)
	assert.True(t, totalEmpty.IsZero())
}

func TestItem_GetTotalActualForYear(t *testing.T) {
	item := createTestItem(t)

	// Add actual amounts for different months in 2024
	err := item.AddActualAmount(2024, 1, mustMoney(t, "95.00", "USD"))
	require.NoError(t, err)
	err = item.AddActualAmount(2024, 2, mustMoney(t, "140.00", "USD"))
	require.NoError(t, err)
	err = item.AddActualAmount(2024, 3, mustMoney(t, "180.00", "USD"))
	require.NoError(t, err)

	// Add actual amount for a different year
	err = item.AddActualAmount(2023, 12, mustMoney(t, "60.00", "USD"))
	require.NoError(t, err)

	total, err := item.GetTotalActualForYear(2024)
	require.NoError(t, err)

	expectedTotal := mustMoney(t, "415.00", "USD") // 95 + 140 + 180
	assert.Equal(t, expectedTotal, total)

	// Test different year
	total2023, err := item.GetTotalActualForYear(2023)
	require.NoError(t, err)

	expectedTotal2023 := mustMoney(t, "60.00", "USD")
	assert.Equal(t, expectedTotal2023, total2023)

	// Test year with no actuals
	totalEmpty, err := item.GetTotalActualForYear(2025)
	require.NoError(t, err)
	assert.True(t, totalEmpty.IsZero())
}

// Helper functions

func createTestItem(t *testing.T) *Item {
	t.Helper()

	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	item, err := NewItem(ledgerID, "Test Item", "Test description", ItemTypeExpense, "USD")
	require.NoError(t, err)

	return item
}

func mustMoney(t *testing.T, amount, currency string) money.Money {
	t.Helper()

	m, err := money.NewMoney(amount, money.Currency(currency))
	require.NoError(t, err)
	return m
}