package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

func TestNewBudgetTracking(t *testing.T) {
	targetAmount := mustMoney(t, "500.00", "USD")

	tests := []struct {
		name         string
		year         int
		month        int
		targetAmount money.Money
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid budget tracking",
			year:         2024,
			month:        6,
			targetAmount: targetAmount,
			wantErr:      false,
		},
		{
			name:         "invalid year - too low",
			year:         1800,
			month:        6,
			targetAmount: targetAmount,
			wantErr:      true,
			errContains:  "invalid year",
		},
		{
			name:         "invalid year - too high",
			year:         3500,
			month:        6,
			targetAmount: targetAmount,
			wantErr:      true,
			errContains:  "invalid year",
		},
		{
			name:         "invalid month - too low",
			year:         2024,
			month:        0,
			targetAmount: targetAmount,
			wantErr:      true,
			errContains:  "invalid month",
		},
		{
			name:         "invalid month - too high",
			year:         2024,
			month:        13,
			targetAmount: targetAmount,
			wantErr:      true,
			errContains:  "invalid month",
		},
		{
			name:         "edge case - minimum valid year",
			year:         1900,
			month:        1,
			targetAmount: targetAmount,
			wantErr:      false,
		},
		{
			name:         "edge case - maximum valid year",
			year:         3000,
			month:        12,
			targetAmount: targetAmount,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt, err := NewBudgetTracking(tt.year, tt.month, tt.targetAmount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, bt)
			} else {
				require.NoError(t, err)
				require.NotNil(t, bt)

				assert.Equal(t, tt.year, bt.Year)
				assert.Equal(t, tt.month, bt.Month)
				assert.Equal(t, tt.targetAmount, bt.TargetAmount)
				assert.Equal(t, tt.targetAmount, bt.BudgetedAmount)
				assert.True(t, bt.ActualAmount.IsZero())
				assert.Equal(t, tt.targetAmount.Currency, bt.ActualAmount.Currency)
				assert.False(t, bt.UpdatedAt.IsZero())
			}
		})
	}
}

func TestReconstructBudgetTracking(t *testing.T) {
	targetAmount := mustMoney(t, "500.00", "USD")
	budgetedAmount := mustMoney(t, "450.00", "USD")
	actualAmount := mustMoney(t, "275.50", "USD")
	updatedAt := time.Now()

	bt := ReconstructBudgetTracking(
		2024,
		7,
		targetAmount,
		budgetedAmount,
		actualAmount,
		updatedAt,
	)

	assert.Equal(t, 2024, bt.Year)
	assert.Equal(t, 7, bt.Month)
	assert.Equal(t, targetAmount, bt.TargetAmount)
	assert.Equal(t, budgetedAmount, bt.BudgetedAmount)
	assert.Equal(t, actualAmount, bt.ActualAmount)
	assert.Equal(t, updatedAt, bt.UpdatedAt)
}

func TestBudgetTracking_UpdateBudgetedAmount(t *testing.T) {
	bt := createTestBudgetTracking(t)
	originalUpdatedAt := bt.UpdatedAt

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid budgeted amount update",
			amount:  mustMoney(t, "400.00", "USD"),
			wantErr: false,
		},
		{
			name:        "currency mismatch",
			amount:      mustMoney(t, "400.00", "EUR"),
			wantErr:     true,
			errContains: "currency mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bt.UpdateBudgetedAmount(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.amount, bt.BudgetedAmount)
				assert.True(t, bt.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestBudgetTracking_AddActualAmount(t *testing.T) {
	bt := createTestBudgetTracking(t)
	originalUpdatedAt := bt.UpdatedAt
	originalActual := bt.ActualAmount

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "add positive amount",
			amount:  mustMoney(t, "100.00", "USD"),
			wantErr: false,
		},
		{
			name:        "currency mismatch",
			amount:      mustMoney(t, "50.00", "EUR"),
			wantErr:     true,
			errContains: "currency mismatch",
		},
		{
			name:    "add negative amount",
			amount:  mustMoney(t, "-25.00", "USD"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeActual := bt.ActualAmount
			err := bt.AddActualAmount(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, beforeActual, bt.ActualAmount)
			} else {
				require.NoError(t, err)
				expectedActual, _ := beforeActual.Add(tt.amount)
				assert.Equal(t, expectedActual, bt.ActualAmount)
				assert.True(t, bt.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}

	// Verify the cumulative effect
	expectedFinalActual, _ := originalActual.Add(mustMoney(t, "100.00", "USD"))
	expectedFinalActual, _ = expectedFinalActual.Add(mustMoney(t, "-25.00", "USD"))
	assert.Equal(t, expectedFinalActual, bt.ActualAmount)
}

func TestBudgetTracking_SubtractActualAmount(t *testing.T) {
	bt := createTestBudgetTracking(t)
	// Add some initial actual amount
	err := bt.AddActualAmount(mustMoney(t, "200.00", "USD"))
	require.NoError(t, err)

	originalUpdatedAt := bt.UpdatedAt
	originalActual := bt.ActualAmount

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "subtract positive amount",
			amount:  mustMoney(t, "50.00", "USD"),
			wantErr: false,
		},
		{
			name:        "currency mismatch",
			amount:      mustMoney(t, "25.00", "EUR"),
			wantErr:     true,
			errContains: "currency mismatch",
		},
		{
			name:    "subtract negative amount",
			amount:  mustMoney(t, "-30.00", "USD"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeActual := bt.ActualAmount
			err := bt.SubtractActualAmount(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, beforeActual, bt.ActualAmount)
			} else {
				require.NoError(t, err)
				expectedActual, _ := beforeActual.Subtract(tt.amount)
				assert.Equal(t, expectedActual, bt.ActualAmount)
				assert.True(t, bt.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}

	// Verify the cumulative effect
	expectedFinalActual, _ := originalActual.Subtract(mustMoney(t, "50.00", "USD"))
	expectedFinalActual, _ = expectedFinalActual.Subtract(mustMoney(t, "-30.00", "USD"))
	assert.Equal(t, expectedFinalActual, bt.ActualAmount)
}

func TestBudgetTracking_GetVariance(t *testing.T) {
	bt := createTestBudgetTracking(t)
	bt.BudgetedAmount = mustMoney(t, "500.00", "USD")
	bt.ActualAmount = mustMoney(t, "350.00", "USD")

	variance, err := bt.GetVariance()
	require.NoError(t, err)

	expectedVariance := mustMoney(t, "150.00", "USD") // 500 - 350
	assert.Equal(t, expectedVariance, variance)

	// Test negative variance (over budget)
	bt.ActualAmount = mustMoney(t, "600.00", "USD")
	variance, err = bt.GetVariance()
	require.NoError(t, err)

	expectedNegativeVariance := mustMoney(t, "-100.00", "USD") // 500 - 600
	assert.Equal(t, expectedNegativeVariance, variance)
}

func TestBudgetTracking_GetTargetVariance(t *testing.T) {
	bt := createTestBudgetTracking(t)
	bt.TargetAmount = mustMoney(t, "400.00", "USD")
	bt.ActualAmount = mustMoney(t, "300.00", "USD")

	variance, err := bt.GetTargetVariance()
	require.NoError(t, err)

	expectedVariance := mustMoney(t, "100.00", "USD") // 400 - 300
	assert.Equal(t, expectedVariance, variance)

	// Test negative variance (over target)
	bt.ActualAmount = mustMoney(t, "450.00", "USD")
	variance, err = bt.GetTargetVariance()
	require.NoError(t, err)

	expectedNegativeVariance := mustMoney(t, "-50.00", "USD") // 400 - 450
	assert.Equal(t, expectedNegativeVariance, variance)
}

func TestBudgetTracking_IsOverBudget(t *testing.T) {
	bt := createTestBudgetTracking(t)

	tests := []struct {
		name           string
		budgetedAmount money.Money
		actualAmount   money.Money
		expected       bool
	}{
		{
			name:           "under budget",
			budgetedAmount: mustMoney(t, "500.00", "USD"),
			actualAmount:   mustMoney(t, "400.00", "USD"),
			expected:       false,
		},
		{
			name:           "exactly on budget",
			budgetedAmount: mustMoney(t, "500.00", "USD"),
			actualAmount:   mustMoney(t, "500.00", "USD"),
			expected:       false,
		},
		{
			name:           "over budget",
			budgetedAmount: mustMoney(t, "500.00", "USD"),
			actualAmount:   mustMoney(t, "600.00", "USD"),
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt.BudgetedAmount = tt.budgetedAmount
			bt.ActualAmount = tt.actualAmount

			result := bt.IsOverBudget()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBudgetTracking_IsOverTarget(t *testing.T) {
	bt := createTestBudgetTracking(t)

	tests := []struct {
		name         string
		targetAmount money.Money
		actualAmount money.Money
		expected     bool
	}{
		{
			name:         "under target",
			targetAmount: mustMoney(t, "400.00", "USD"),
			actualAmount: mustMoney(t, "350.00", "USD"),
			expected:     false,
		},
		{
			name:         "exactly on target",
			targetAmount: mustMoney(t, "400.00", "USD"),
			actualAmount: mustMoney(t, "400.00", "USD"),
			expected:     false,
		},
		{
			name:         "over target",
			targetAmount: mustMoney(t, "400.00", "USD"),
			actualAmount: mustMoney(t, "450.00", "USD"),
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt.TargetAmount = tt.targetAmount
			bt.ActualAmount = tt.actualAmount

			result := bt.IsOverTarget()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBudgetTracking_GetBudgetUtilization(t *testing.T) {
	bt := createTestBudgetTracking(t)

	tests := []struct {
		name           string
		budgetedAmount money.Money
		actualAmount   money.Money
		expected       float64
	}{
		{
			name:           "50% utilization",
			budgetedAmount: mustMoney(t, "100.00", "USD"),
			actualAmount:   mustMoney(t, "50.00", "USD"),
			expected:       0.5,
		},
		{
			name:           "100% utilization",
			budgetedAmount: mustMoney(t, "200.00", "USD"),
			actualAmount:   mustMoney(t, "200.00", "USD"),
			expected:       1.0,
		},
		{
			name:           "150% utilization (over budget)",
			budgetedAmount: mustMoney(t, "100.00", "USD"),
			actualAmount:   mustMoney(t, "150.00", "USD"),
			expected:       1.5,
		},
		{
			name:           "zero budget",
			budgetedAmount: mustMoney(t, "0.00", "USD"),
			actualAmount:   mustMoney(t, "50.00", "USD"),
			expected:       0.0,
		},
		{
			name:           "zero actual",
			budgetedAmount: mustMoney(t, "100.00", "USD"),
			actualAmount:   mustMoney(t, "0.00", "USD"),
			expected:       0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt.BudgetedAmount = tt.budgetedAmount
			bt.ActualAmount = tt.actualAmount

			result := bt.GetBudgetUtilization()
			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

func TestBudgetTracking_GetTargetUtilization(t *testing.T) {
	bt := createTestBudgetTracking(t)

	tests := []struct {
		name         string
		targetAmount money.Money
		actualAmount money.Money
		expected     float64
	}{
		{
			name:         "75% of target",
			targetAmount: mustMoney(t, "400.00", "USD"),
			actualAmount: mustMoney(t, "300.00", "USD"),
			expected:     0.75,
		},
		{
			name:         "exactly on target",
			targetAmount: mustMoney(t, "500.00", "USD"),
			actualAmount: mustMoney(t, "500.00", "USD"),
			expected:     1.0,
		},
		{
			name:         "120% of target (over target)",
			targetAmount: mustMoney(t, "100.00", "USD"),
			actualAmount: mustMoney(t, "120.00", "USD"),
			expected:     1.2,
		},
		{
			name:         "zero target",
			targetAmount: mustMoney(t, "0.00", "USD"),
			actualAmount: mustMoney(t, "50.00", "USD"),
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt.TargetAmount = tt.targetAmount
			bt.ActualAmount = tt.actualAmount

			result := bt.GetTargetUtilization()
			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

func TestBudgetTracking_GetRemainingBudget(t *testing.T) {
	bt := createTestBudgetTracking(t)

	tests := []struct {
		name           string
		budgetedAmount money.Money
		actualAmount   money.Money
		expectedAmount money.Money
		expectZero     bool
	}{
		{
			name:           "remaining budget available",
			budgetedAmount: mustMoney(t, "500.00", "USD"),
			actualAmount:   mustMoney(t, "300.00", "USD"),
			expectedAmount: mustMoney(t, "200.00", "USD"),
			expectZero:     false,
		},
		{
			name:           "exactly on budget",
			budgetedAmount: mustMoney(t, "400.00", "USD"),
			actualAmount:   mustMoney(t, "400.00", "USD"),
			expectedAmount: mustMoney(t, "0.00", "USD"),
			expectZero:     true,
		},
		{
			name:           "over budget (should return zero)",
			budgetedAmount: mustMoney(t, "300.00", "USD"),
			actualAmount:   mustMoney(t, "350.00", "USD"),
			expectedAmount: mustMoney(t, "0.00", "USD"),
			expectZero:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt.BudgetedAmount = tt.budgetedAmount
			bt.ActualAmount = tt.actualAmount

			remaining, err := bt.GetRemainingBudget()
			require.NoError(t, err)

			assert.True(t, remaining.Equals(tt.expectedAmount))
			assert.Equal(t, tt.expectZero, remaining.IsZero())
		})
	}
}

// Helper functions

func createTestBudgetTracking(t *testing.T) *BudgetTracking {
	t.Helper()

	targetAmount := mustMoney(t, "500.00", "USD")
	bt, err := NewBudgetTracking(2024, 6, targetAmount)
	require.NoError(t, err)

	return bt
}