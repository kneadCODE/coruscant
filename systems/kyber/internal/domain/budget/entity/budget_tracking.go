package entity

import (
	"fmt"
	"time"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

// BudgetTracking represents the budget tracking for an item in a specific month
type BudgetTracking struct {
	Year           int
	Month          int
	TargetAmount   money.Money // Original planned amount
	BudgetedAmount money.Money // Current approved budget (can be adjusted)
	ActualAmount   money.Money // Actual amount from transactions
	UpdatedAt      time.Time
}

// NewBudgetTracking creates a new BudgetTracking for a specific month
func NewBudgetTracking(year, month int, targetAmount money.Money) (*BudgetTracking, error) {
	if year < 1900 || year > 3000 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}

	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d", month)
	}

	// Initialize budgeted amount to target amount
	// Actual amount starts at zero
	actualAmount, err := money.Zero(targetAmount.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize actual amount: %w", err)
	}

	return &BudgetTracking{
		Year:           year,
		Month:          month,
		TargetAmount:   targetAmount,
		BudgetedAmount: targetAmount, // Start with target as budgeted
		ActualAmount:   actualAmount,
		UpdatedAt:      time.Now(),
	}, nil
}

// ReconstructBudgetTracking reconstructs BudgetTracking from stored data
func ReconstructBudgetTracking(
	year, month int,
	targetAmount, budgetedAmount, actualAmount money.Money,
	updatedAt time.Time,
) *BudgetTracking {
	return &BudgetTracking{
		Year:           year,
		Month:          month,
		TargetAmount:   targetAmount,
		BudgetedAmount: budgetedAmount,
		ActualAmount:   actualAmount,
		UpdatedAt:      updatedAt,
	}
}

// UpdateBudgetedAmount updates the budgeted amount for this month
func (bt *BudgetTracking) UpdateBudgetedAmount(amount money.Money) error {
	if amount.Currency != bt.TargetAmount.Currency {
		return fmt.Errorf("currency mismatch: expected %s, got %s", bt.TargetAmount.Currency, amount.Currency)
	}

	bt.BudgetedAmount = amount
	bt.UpdatedAt = time.Now()
	return nil
}

// AddActualAmount adds to the actual amount (from transactions)
func (bt *BudgetTracking) AddActualAmount(amount money.Money) error {
	if amount.Currency != bt.ActualAmount.Currency {
		return fmt.Errorf("currency mismatch: expected %s, got %s", bt.ActualAmount.Currency, amount.Currency)
	}

	newActual, err := bt.ActualAmount.Add(amount)
	if err != nil {
		return fmt.Errorf("failed to add actual amount: %w", err)
	}

	bt.ActualAmount = newActual
	bt.UpdatedAt = time.Now()
	return nil
}

// SubtractActualAmount subtracts from the actual amount (transaction reversal)
func (bt *BudgetTracking) SubtractActualAmount(amount money.Money) error {
	if amount.Currency != bt.ActualAmount.Currency {
		return fmt.Errorf("currency mismatch: expected %s, got %s", bt.ActualAmount.Currency, amount.Currency)
	}

	newActual, err := bt.ActualAmount.Subtract(amount)
	if err != nil {
		return fmt.Errorf("failed to subtract actual amount: %w", err)
	}

	bt.ActualAmount = newActual
	bt.UpdatedAt = time.Now()
	return nil
}

// GetVariance returns the variance between budgeted and actual amounts
func (bt *BudgetTracking) GetVariance() (money.Money, error) {
	return bt.BudgetedAmount.Subtract(bt.ActualAmount)
}

// GetTargetVariance returns the variance between target and actual amounts
func (bt *BudgetTracking) GetTargetVariance() (money.Money, error) {
	return bt.TargetAmount.Subtract(bt.ActualAmount)
}

// IsOverBudget checks if actual amount exceeds budgeted amount
func (bt *BudgetTracking) IsOverBudget() bool {
	variance, err := bt.GetVariance()
	if err != nil {
		return false
	}
	return variance.IsNegative()
}

// IsOverTarget checks if actual amount exceeds target amount
func (bt *BudgetTracking) IsOverTarget() bool {
	variance, err := bt.GetTargetVariance()
	if err != nil {
		return false
	}
	return variance.IsNegative()
}

// GetBudgetUtilization returns the percentage of budget used (0.0 to 1.0+)
func (bt *BudgetTracking) GetBudgetUtilization() float64 {
	if bt.BudgetedAmount.IsZero() {
		return 0.0
	}
	return bt.ActualAmount.Float64() / bt.BudgetedAmount.Float64()
}

// GetTargetUtilization returns the percentage of target used (0.0 to 1.0+)
func (bt *BudgetTracking) GetTargetUtilization() float64 {
	if bt.TargetAmount.IsZero() {
		return 0.0
	}
	return bt.ActualAmount.Float64() / bt.TargetAmount.Float64()
}

// GetRemainingBudget returns the remaining budget amount
func (bt *BudgetTracking) GetRemainingBudget() (money.Money, error) {
	remaining, err := bt.BudgetedAmount.Subtract(bt.ActualAmount)
	if err != nil {
		return money.Money{}, err
	}

	// Return zero if over budget (no negative remaining)
	if remaining.IsNegative() {
		return money.Zero(bt.BudgetedAmount.Currency)
	}

	return remaining, nil
}
