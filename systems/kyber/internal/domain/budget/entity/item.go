package entity

import (
	"fmt"
	"time"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

// Item represents a budget item (income, expense, or transfer) within a ledger
type Item struct {
	ID          ItemID
	LedgerID    entity.LedgerID
	Name        string
	Description string
	Type        ItemType
	// Category       ItemCategory
	Currency       money.Currency
	MonthlyBudgets map[string]*BudgetTracking // Key: "YYYY-MM"
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewItem creates a new Item
func NewItem(
	ledgerID entity.LedgerID,
	name, description string,
	itemType ItemType,
	// category ItemCategory,
	currency money.Currency,
) (*Item, error) {
	if !ledgerID.IsValid() {
		return nil, fmt.Errorf("ledger ID cannot be empty")
	}

	if name == "" {
		return nil, fmt.Errorf("item name cannot be empty")
	}

	if currency == "" {
		return nil, fmt.Errorf("item currency cannot be empty")
	}

	// Validate category matches type
	// if category.GetItemType() != itemType {
	// 	return nil, fmt.Errorf("category %s does not match item type %s", category, itemType)
	// }

	id, err := NewItemID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate item ID: %w", err)
	}

	now := time.Now()

	return &Item{
		ID:          id,
		LedgerID:    ledgerID,
		Name:        name,
		Description: description,
		Type:        itemType,
		// Category:       category,
		Currency:       currency,
		MonthlyBudgets: make(map[string]*BudgetTracking),
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// ReconstructItem reconstructs an Item from stored data
func ReconstructItem(
	id ItemID,
	ledgerID entity.LedgerID,
	name, description string,
	itemType ItemType,
	// category ItemCategory,
	currency money.Currency,
	monthlyBudgets map[string]*BudgetTracking,
	isActive bool,
	createdAt, updatedAt time.Time,
) *Item {
	if monthlyBudgets == nil {
		monthlyBudgets = make(map[string]*BudgetTracking)
	}

	return &Item{
		ID:          id,
		LedgerID:    ledgerID,
		Name:        name,
		Description: description,
		Type:        itemType,
		// Category:       category,
		Currency:       currency,
		MonthlyBudgets: monthlyBudgets,
		IsActive:       isActive,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

// UpdateInfo updates the item's basic information
func (i *Item) UpdateInfo(name, description string) error {
	if name == "" {
		return fmt.Errorf("item name cannot be empty")
	}

	i.Name = name
	i.Description = description
	i.UpdatedAt = time.Now()
	return nil
}

// // UpdateCategory updates the item's category (must match type)
// func (i *Item) UpdateCategory(category ItemCategory) error {
// 	if category.GetItemType() != i.Type {
// 		return fmt.Errorf("category %s does not match item type %s", category, i.Type)
// 	}

// 	i.Category = category
// 	i.UpdatedAt = time.Now()
// 	return nil
// }

// SetMonthlyTarget sets the target amount for a specific month
func (i *Item) SetMonthlyTarget(year, month int, targetAmount money.Money) error {
	if targetAmount.Currency != i.Currency {
		return fmt.Errorf("currency mismatch: item uses %s, target uses %s", i.Currency, targetAmount.Currency)
	}

	monthKey := fmt.Sprintf("%04d-%02d", year, month)

	if existing, exists := i.MonthlyBudgets[monthKey]; exists {
		// Update existing target and reset budgeted to new target
		existing.TargetAmount = targetAmount
		existing.BudgetedAmount = targetAmount
		existing.UpdatedAt = time.Now()
	} else {
		// Create new budget tracking
		budgetTracking, err := NewBudgetTracking(year, month, targetAmount)
		if err != nil {
			return fmt.Errorf("failed to create budget tracking: %w", err)
		}
		i.MonthlyBudgets[monthKey] = budgetTracking
	}

	i.UpdatedAt = time.Now()
	return nil
}

// UpdateMonthlyBudget updates the budgeted amount for a specific month
func (i *Item) UpdateMonthlyBudget(year, month int, budgetedAmount money.Money) error {
	monthKey := fmt.Sprintf("%04d-%02d", year, month)

	budgetTracking, exists := i.MonthlyBudgets[monthKey]
	if !exists {
		return fmt.Errorf("no budget tracking found for %s", monthKey)
	}

	return budgetTracking.UpdateBudgetedAmount(budgetedAmount)
}

// AddActualAmount adds actual spending/income to a specific month
func (i *Item) AddActualAmount(year, month int, amount money.Money) error {
	monthKey := fmt.Sprintf("%04d-%02d", year, month)

	budgetTracking, exists := i.MonthlyBudgets[monthKey]
	if !exists {
		// Create budget tracking with zero target if it doesn't exist
		zeroTarget, err := money.Zero(i.Currency)
		if err != nil {
			return fmt.Errorf("failed to create zero target: %w", err)
		}

		budgetTracking, err = NewBudgetTracking(year, month, zeroTarget)
		if err != nil {
			return fmt.Errorf("failed to create budget tracking: %w", err)
		}
		i.MonthlyBudgets[monthKey] = budgetTracking
	}

	err := budgetTracking.AddActualAmount(amount)
	if err != nil {
		return err
	}

	i.UpdatedAt = time.Now()
	return nil
}

// GetMonthlyBudget returns the budget tracking for a specific month
func (i *Item) GetMonthlyBudget(year, month int) *BudgetTracking {
	monthKey := fmt.Sprintf("%04d-%02d", year, month)
	return i.MonthlyBudgets[monthKey]
}

// GetCurrentMonthBudget returns the budget tracking for the current month
func (i *Item) GetCurrentMonthBudget() *BudgetTracking {
	now := time.Now()
	return i.GetMonthlyBudget(now.Year(), int(now.Month()))
}

// Activate activates the item
func (i *Item) Activate() {
	i.IsActive = true
	i.UpdatedAt = time.Now()
}

// Deactivate deactivates the item
func (i *Item) Deactivate() {
	i.IsActive = false
	i.UpdatedAt = time.Now()
}

// GetTotalBudgetedForYear returns the total budgeted amount for a specific year
func (i *Item) GetTotalBudgetedForYear(year int) (money.Money, error) {
	total, err := money.Zero(i.Currency)
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to initialize total: %w", err)
	}

	for monthKey, budget := range i.MonthlyBudgets {
		if len(monthKey) >= 4 && monthKey[:4] == fmt.Sprintf("%04d", year) {
			total, err = total.Add(budget.BudgetedAmount)
			if err != nil {
				return money.Money{}, fmt.Errorf("failed to add budget amount: %w", err)
			}
		}
	}

	return total, nil
}

// GetTotalActualForYear returns the total actual amount for a specific year
func (i *Item) GetTotalActualForYear(year int) (money.Money, error) {
	total, err := money.Zero(i.Currency)
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to initialize total: %w", err)
	}

	for monthKey, budget := range i.MonthlyBudgets {
		if len(monthKey) >= 4 && monthKey[:4] == fmt.Sprintf("%04d", year) {
			total, err = total.Add(budget.ActualAmount)
			if err != nil {
				return money.Money{}, fmt.Errorf("failed to add actual amount: %w", err)
			}
		}
	}

	return total, nil
}
