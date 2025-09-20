package entity

import (
	"fmt"
	"time"

	"github.com/kneadCODE/coruscant/shared/golib/optional"
	budgetEntity "github.com/kneadCODE/coruscant/systems/kyber/internal/domain/budget/entity"
	counterpartyEntity "github.com/kneadCODE/coruscant/systems/kyber/internal/domain/counterparty/entity"
	ledgerEntity "github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

// Transaction represents a financial transaction within a ledger
type Transaction struct {
	ID              TransactionID
	LedgerID        ledgerEntity.LedgerID
	AccountID       AccountID
	ItemID          budgetEntity.ItemID
	CounterpartyID  optional.Option[counterpartyEntity.CounterpartyID] // Optional - who the transaction is with
	Amount          money.Money
	Description     string
	Notes           string
	TransactionDate time.Time // When the transaction actually occurred
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewTransaction creates a new Transaction
func NewTransaction(
	ledgerID ledgerEntity.LedgerID,
	accountID AccountID,
	itemID budgetEntity.ItemID,
	amount money.Money,
	description string,
	transactionDate time.Time,
) (*Transaction, error) {
	if !ledgerID.IsValid() {
		return nil, fmt.Errorf("ledger ID is invalid")
	}

	if !accountID.IsValid() {
		return nil, fmt.Errorf("account ID is invalid")
	}

	if !itemID.IsValid() {
		return nil, fmt.Errorf("item ID is invalid")
	}

	if amount.IsZero() {
		return nil, fmt.Errorf("transaction amount cannot be zero")
	}

	if description == "" {
		return nil, fmt.Errorf("transaction description cannot be empty")
	}

	if transactionDate.IsZero() {
		transactionDate = time.Now()
	}

	id, err := NewTransactionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ID: %w", err)
	}

	now := time.Now()

	return &Transaction{
		ID:              id,
		LedgerID:        ledgerID,
		AccountID:       accountID,
		ItemID:          itemID,
		CounterpartyID:  optional.None[counterpartyEntity.CounterpartyID](),
		Amount:          amount,
		Description:     description,
		TransactionDate: transactionDate,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// NewTransactionWithCounterparty creates a new Transaction with a counterparty
func NewTransactionWithCounterparty(
	ledgerID ledgerEntity.LedgerID,
	accountID AccountID,
	itemID budgetEntity.ItemID,
	counterpartyID counterpartyEntity.CounterpartyID,
	amount money.Money,
	description string,
	transactionDate time.Time,
) (*Transaction, error) {
	tx, err := NewTransaction(ledgerID, accountID, itemID, amount, description, transactionDate)
	if err != nil {
		return nil, err
	}
	tx.SetCounterparty(counterpartyID)
	return tx, nil
}

// ReconstructTransaction reconstructs a Transaction from stored data
func ReconstructTransaction(
	id TransactionID,
	ledgerID ledgerEntity.LedgerID,
	accountID AccountID,
	itemID budgetEntity.ItemID,
	counterpartyID optional.Option[counterpartyEntity.CounterpartyID],
	amount money.Money,
	description, notes string,
	transactionDate, createdAt, updatedAt time.Time,
) *Transaction {
	return &Transaction{
		ID:              id,
		LedgerID:        ledgerID,
		AccountID:       accountID,
		ItemID:          itemID,
		CounterpartyID:  counterpartyID,
		Amount:          amount,
		Description:     description,
		Notes:           notes,
		TransactionDate: transactionDate,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// UpdateInfo updates the transaction's basic information
func (t *Transaction) UpdateInfo(description, notes string) error {
	if description == "" {
		return fmt.Errorf("transaction description cannot be empty")
	}

	t.Description = description
	t.Notes = notes
	t.UpdatedAt = time.Now()
	return nil
}

// UpdateAmount updates the transaction amount
func (t *Transaction) UpdateAmount(amount money.Money) error {
	if amount.IsZero() {
		return fmt.Errorf("transaction amount cannot be zero")
	}

	if amount.Currency != t.Amount.Currency {
		return fmt.Errorf("currency mismatch: transaction uses %s, new amount uses %s", t.Amount.Currency, amount.Currency)
	}

	t.Amount = amount
	t.UpdatedAt = time.Now()
	return nil
}

// UpdateTransactionDate updates when the transaction occurred
func (t *Transaction) UpdateTransactionDate(transactionDate time.Time) {
	t.TransactionDate = transactionDate
	t.UpdatedAt = time.Now()
}

// UpdateItem udpdates the budget item
func (t *Transaction) UpdateItem(itemID budgetEntity.ItemID) {
	t.ItemID = itemID
	t.UpdatedAt = time.Now()
}

// SetCounterparty associates the transaction with a counterparty
func (t *Transaction) SetCounterparty(counterpartyID counterpartyEntity.CounterpartyID) {
	t.CounterpartyID = optional.Some(counterpartyID)
	t.UpdatedAt = time.Now()
}

// RemoveCounterparty removes the association with a counterparty
func (t *Transaction) RemoveCounterparty() {
	t.CounterpartyID = optional.None[counterpartyEntity.CounterpartyID]()
	t.UpdatedAt = time.Now()
}

// HasCounterparty checks if the transaction is associated with a counterparty
func (t *Transaction) HasCounterparty() bool {
	return t.CounterpartyID.IsSome()
}

// GetCounterparty returns the counterparty ID if present
func (t *Transaction) GetCounterparty() (counterpartyEntity.CounterpartyID, bool) {
	if t.CounterpartyID.IsSome() {
		return t.CounterpartyID.Unwrap(), true
	}
	return counterpartyEntity.CounterpartyID{}, false
}

// GetCounterpartyID returns the counterparty ID as an Option
func (t *Transaction) GetCounterpartyID() optional.Option[counterpartyEntity.CounterpartyID] {
	return t.CounterpartyID
}

// IsDebit checks if the transaction is a debit (negative amount)
func (t *Transaction) IsDebit() bool {
	return t.Amount.IsNegative()
}

// IsCredit checks if the transaction is a credit (positive amount)
func (t *Transaction) IsCredit() bool {
	return t.Amount.IsPositive()
}

// GetAbsoluteAmount returns the absolute value of the transaction amount
func (t *Transaction) GetAbsoluteAmount() money.Money {
	return t.Amount.Abs()
}
