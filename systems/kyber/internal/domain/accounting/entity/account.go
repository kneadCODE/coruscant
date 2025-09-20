package entity

import (
	"fmt"
	"time"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

// Account represents a financial account within a ledger
type Account struct {
	ID          AccountID
	LedgerID    entity.LedgerID
	Name        string
	Description string
	Type        AccountType
	Currency    money.Currency
	Balance     money.Money // Current balance
	Status      AccountStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewAccount creates a new Account
func NewAccount(
	ledgerID entity.LedgerID,
	name, description string,
	accountType AccountType,
	currency money.Currency,
) (*Account, error) {
	if !ledgerID.IsValid() {
		return nil, fmt.Errorf("ledger ID is invalid")
	}

	if name == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	if currency == "" {
		return nil, fmt.Errorf("account currency cannot be empty")
	}

	// Initialize with zero balance
	balance, err := money.Zero(currency)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize balance: %w", err)
	}

	id, err := NewAccountID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate account ID: %w", err)
	}

	now := time.Now()

	return &Account{
		ID:          id,
		LedgerID:    ledgerID,
		Name:        name,
		Description: description,
		Type:        accountType,
		Currency:    currency,
		Balance:     balance,
		Status:      AccountStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ReconstructAccount reconstructs an Account from stored data
func ReconstructAccount(
	id AccountID,
	ledgerID entity.LedgerID,
	name, description string,
	accountType AccountType,
	currency money.Currency,
	balance money.Money,
	status AccountStatus,
	createdAt, updatedAt time.Time,
) *Account {
	return &Account{
		ID:          id,
		LedgerID:    ledgerID,
		Name:        name,
		Description: description,
		Type:        accountType,
		Currency:    currency,
		Balance:     balance,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// UpdateInfo updates the account's basic information
func (a *Account) UpdateInfo(name, description string) error {
	if name == "" {
		return fmt.Errorf("account name cannot be empty")
	}

	a.Name = name
	a.Description = description
	a.UpdatedAt = time.Now()
	return nil
}

// Credit adds money to the account (increases balance)
func (a *Account) Credit(amount money.Money) error {
	if amount.Currency != a.Currency {
		return fmt.Errorf("currency mismatch: account uses %s, transaction uses %s", a.Currency, amount.Currency)
	}

	if amount.IsNegative() {
		return fmt.Errorf("credit amount cannot be negative")
	}

	newBalance, err := a.Balance.Add(amount)
	if err != nil {
		return fmt.Errorf("failed to credit account: %w", err)
	}

	a.Balance = newBalance
	a.UpdatedAt = time.Now()
	return nil
}

// Debit removes money from the account (decreases balance)
func (a *Account) Debit(amount money.Money) error {
	if amount.Currency != a.Currency {
		return fmt.Errorf("currency mismatch: account uses %s, transaction uses %s", a.Currency, amount.Currency)
	}

	if amount.IsNegative() {
		return fmt.Errorf("debit amount cannot be negative")
	}

	newBalance, err := a.Balance.Subtract(amount)
	if err != nil {
		return fmt.Errorf("failed to debit account: %w", err)
	}

	a.Balance = newBalance
	a.UpdatedAt = time.Now()
	return nil
}

// Activate activates the account
func (a *Account) Activate() {
	a.Status = AccountStatusActive
	a.UpdatedAt = time.Now()
}

// Archive archives the account
func (a *Account) Archive() {
	a.Status = AccountStatusArchived
	a.UpdatedAt = time.Now()
}

// CanDebit checks if the account can be debited by the specified amount
// For liability accounts (credit cards, loans), this allows negative balances
func (a *Account) CanDebit(amount money.Money) bool {
	if amount.Currency != a.Currency {
		return false
	}

	if amount.IsNegative() {
		return false
	}

	// For liability accounts, negative balances are normal
	if a.Type.IsLiability() {
		return true
	}

	// For asset accounts, check if sufficient balance exists
	newBalance, err := a.Balance.Subtract(amount)
	if err != nil {
		return false
	}

	return !newBalance.IsNegative()
}

// HasSufficientBalance checks if the account has sufficient balance for a debit
func (a *Account) HasSufficientBalance(amount money.Money) bool {
	return a.CanDebit(amount)
}

// GetBalanceFloat64 returns the balance as a float64 for display purposes
func (a *Account) GetBalanceFloat64() float64 {
	return a.Balance.Float64()
}

// DebitBalance reduces the account balance by the specified amount
func (a *Account) DebitBalance(amount money.Money) error {
	if !a.CanDebit(amount) {
		return fmt.Errorf("insufficient balance: %s available, %s required", a.Balance, amount)
	}

	newBalance, err := a.Balance.Subtract(amount)
	if err != nil {
		return fmt.Errorf("failed to debit balance: %w", err)
	}

	a.Balance = newBalance
	a.UpdatedAt = time.Now()
	return nil
}

// CreditBalance increases the account balance by the specified amount
func (a *Account) CreditBalance(amount money.Money) error {
	newBalance, err := a.Balance.Add(amount)
	if err != nil {
		return fmt.Errorf("failed to credit balance: %w", err)
	}

	a.Balance = newBalance
	a.UpdatedAt = time.Now()
	return nil
}
