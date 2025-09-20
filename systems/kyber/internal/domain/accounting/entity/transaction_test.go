package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/shared/golib/optional"
	budgetEntity "github.com/kneadCODE/coruscant/systems/kyber/internal/domain/budget/entity"
	counterpartyEntity "github.com/kneadCODE/coruscant/systems/kyber/internal/domain/counterparty/entity"
	ledgerEntity "github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

func TestNewTransaction(t *testing.T) {
	ledgerID, err := ledgerEntity.NewLedgerID()
	require.NoError(t, err)

	accountID, err := NewAccountID()
	require.NoError(t, err)

	itemID, err := budgetEntity.NewItemID()
	require.NoError(t, err)

	amount := mustMoney(t, "100.50", "USD")
	transactionDate := time.Now().Add(-time.Hour)

	tests := []struct {
		name            string
		ledgerID        ledgerEntity.LedgerID
		accountID       AccountID
		itemID          budgetEntity.ItemID
		amount          money.Money
		description     string
		transactionDate time.Time
		wantErr         bool
		errContains     string
	}{
		{
			name:            "valid transaction",
			ledgerID:        ledgerID,
			accountID:       accountID,
			itemID:          itemID,
			amount:          amount,
			description:     "Test transaction",
			transactionDate: transactionDate,
			wantErr:         false,
		},
		{
			name:            "invalid ledger ID",
			ledgerID:        ledgerEntity.LedgerID{},
			accountID:       accountID,
			itemID:          itemID,
			amount:          amount,
			description:     "Test transaction",
			transactionDate: transactionDate,
			wantErr:         true,
			errContains:     "ledger ID is invalid",
		},
		{
			name:            "invalid account ID",
			ledgerID:        ledgerID,
			accountID:       AccountID{},
			itemID:          itemID,
			amount:          amount,
			description:     "Test transaction",
			transactionDate: transactionDate,
			wantErr:         true,
			errContains:     "account ID is invalid",
		},
		{
			name:            "invalid item ID",
			ledgerID:        ledgerID,
			accountID:       accountID,
			itemID:          budgetEntity.ItemID{},
			amount:          amount,
			description:     "Test transaction",
			transactionDate: transactionDate,
			wantErr:         true,
			errContains:     "item ID is invalid",
		},
		{
			name:            "zero amount",
			ledgerID:        ledgerID,
			accountID:       accountID,
			itemID:          itemID,
			amount:          mustMoney(t, "0.00", "USD"),
			description:     "Test transaction",
			transactionDate: transactionDate,
			wantErr:         true,
			errContains:     "transaction amount cannot be zero",
		},
		{
			name:            "empty description",
			ledgerID:        ledgerID,
			accountID:       accountID,
			itemID:          itemID,
			amount:          amount,
			description:     "",
			transactionDate: transactionDate,
			wantErr:         true,
			errContains:     "transaction description cannot be empty",
		},
		{
			name:            "zero transaction date defaults to now",
			ledgerID:        ledgerID,
			accountID:       accountID,
			itemID:          itemID,
			amount:          amount,
			description:     "Test transaction",
			transactionDate: time.Time{},
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewTransaction(
				tt.ledgerID,
				tt.accountID,
				tt.itemID,
				tt.amount,
				tt.description,
				tt.transactionDate,
			)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, transaction)
			} else {
				require.NoError(t, err)
				require.NotNil(t, transaction)

				assert.True(t, transaction.ID.IsValid())
				assert.Equal(t, tt.ledgerID, transaction.LedgerID)
				assert.Equal(t, tt.accountID, transaction.AccountID)
				assert.Equal(t, tt.itemID, transaction.ItemID)
				assert.Equal(t, tt.amount, transaction.Amount)
				assert.Equal(t, tt.description, transaction.Description)
				assert.Empty(t, transaction.Notes)
				assert.False(t, transaction.HasCounterparty())
				assert.False(t, transaction.CreatedAt.IsZero())
				assert.False(t, transaction.UpdatedAt.IsZero())

				if tt.transactionDate.IsZero() {
					// Should default to current time (within reasonable range)
					assert.True(t, time.Since(transaction.TransactionDate) < time.Minute)
				} else {
					assert.Equal(t, tt.transactionDate, transaction.TransactionDate)
				}
			}
		})
	}
}

func TestNewTransactionWithCounterparty(t *testing.T) {
	ledgerID, err := ledgerEntity.NewLedgerID()
	require.NoError(t, err)

	accountID, err := NewAccountID()
	require.NoError(t, err)

	itemID, err := budgetEntity.NewItemID()
	require.NoError(t, err)

	counterpartyID, err := counterpartyEntity.NewCounterpartyID()
	require.NoError(t, err)

	amount := mustMoney(t, "100.50", "USD")
	transactionDate := time.Now().Add(-time.Hour)

	transaction, err := NewTransactionWithCounterparty(
		ledgerID,
		accountID,
		itemID,
		counterpartyID,
		amount,
		"Test transaction with counterparty",
		transactionDate,
	)

	require.NoError(t, err)
	require.NotNil(t, transaction)

	assert.True(t, transaction.HasCounterparty())
	retrievedCounterpartyID, ok := transaction.GetCounterparty()
	assert.True(t, ok)
	assert.Equal(t, counterpartyID, retrievedCounterpartyID)
}

func TestReconstructTransaction(t *testing.T) {
	transactionID, err := NewTransactionID()
	require.NoError(t, err)

	ledgerID, err := ledgerEntity.NewLedgerID()
	require.NoError(t, err)

	accountID, err := NewAccountID()
	require.NoError(t, err)

	itemID, err := budgetEntity.NewItemID()
	require.NoError(t, err)

	counterpartyID, err := counterpartyEntity.NewCounterpartyID()
	require.NoError(t, err)

	amount := mustMoney(t, "250.75", "USD")
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()
	transactionDate := time.Now().Add(-30 * time.Minute)

	transaction := ReconstructTransaction(
		transactionID,
		ledgerID,
		accountID,
		itemID,
		optional.Some(counterpartyID),
		amount,
		"Reconstructed transaction",
		"Some notes",
		transactionDate,
		createdAt,
		updatedAt,
	)

	assert.Equal(t, transactionID, transaction.ID)
	assert.Equal(t, ledgerID, transaction.LedgerID)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, itemID, transaction.ItemID)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, "Reconstructed transaction", transaction.Description)
	assert.Equal(t, "Some notes", transaction.Notes)
	assert.Equal(t, transactionDate, transaction.TransactionDate)
	assert.Equal(t, createdAt, transaction.CreatedAt)
	assert.Equal(t, updatedAt, transaction.UpdatedAt)
	assert.True(t, transaction.HasCounterparty())

	retrievedCounterpartyID, ok := transaction.GetCounterparty()
	assert.True(t, ok)
	assert.Equal(t, counterpartyID, retrievedCounterpartyID)
}

func TestTransaction_UpdateInfo(t *testing.T) {
	transaction := createTestTransaction(t)
	originalUpdatedAt := transaction.UpdatedAt

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		description string
		notes       string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid update",
			description: "Updated description",
			notes:       "Updated notes",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			notes:       "Some notes",
			wantErr:     true,
			errContains: "transaction description cannot be empty",
		},
		{
			name:        "empty notes is allowed",
			description: "Valid description",
			notes:       "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transaction.UpdateInfo(tt.description, tt.notes)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.description, transaction.Description)
				assert.Equal(t, tt.notes, transaction.Notes)
				assert.True(t, transaction.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestTransaction_UpdateAmount(t *testing.T) {
	transaction := createTestTransaction(t)
	originalUpdatedAt := transaction.UpdatedAt

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid amount update",
			amount:  mustMoney(t, "200.00", "USD"),
			wantErr: false,
		},
		{
			name:        "zero amount",
			amount:      mustMoney(t, "0.00", "USD"),
			wantErr:     true,
			errContains: "transaction amount cannot be zero",
		},
		{
			name:        "currency mismatch",
			amount:      mustMoney(t, "100.00", "EUR"),
			wantErr:     true,
			errContains: "currency mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transaction.UpdateAmount(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.amount, transaction.Amount)
				assert.True(t, transaction.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestTransaction_UpdateTransactionDate(t *testing.T) {
	transaction := createTestTransaction(t)
	originalUpdatedAt := transaction.UpdatedAt

	time.Sleep(time.Millisecond)

	newDate := time.Now().Add(-2 * time.Hour)
	transaction.UpdateTransactionDate(newDate)

	assert.Equal(t, newDate, transaction.TransactionDate)
	assert.True(t, transaction.UpdatedAt.After(originalUpdatedAt))
}

func TestTransaction_UpdateItem(t *testing.T) {
	transaction := createTestTransaction(t)
	originalUpdatedAt := transaction.UpdatedAt

	newItemID, err := budgetEntity.NewItemID()
	require.NoError(t, err)

	time.Sleep(time.Millisecond)

	transaction.UpdateItem(newItemID)

	assert.Equal(t, newItemID, transaction.ItemID)
	assert.True(t, transaction.UpdatedAt.After(originalUpdatedAt))
}

func TestTransaction_CounterpartyOperations(t *testing.T) {
	transaction := createTestTransaction(t)

	// Initially should have no counterparty
	assert.False(t, transaction.HasCounterparty())

	counterpartyID, hasCounterparty := transaction.GetCounterparty()
	assert.False(t, hasCounterparty)
	assert.Equal(t, counterpartyEntity.CounterpartyID{}, counterpartyID)

	assert.True(t, transaction.GetCounterpartyID().IsNone())

	// Set counterparty
	newCounterpartyID, err := counterpartyEntity.NewCounterpartyID()
	require.NoError(t, err)

	originalUpdatedAt := transaction.UpdatedAt
	time.Sleep(time.Millisecond)

	transaction.SetCounterparty(newCounterpartyID)

	assert.True(t, transaction.HasCounterparty())
	assert.True(t, transaction.UpdatedAt.After(originalUpdatedAt))

	retrievedCounterpartyID, hasCounterparty := transaction.GetCounterparty()
	assert.True(t, hasCounterparty)
	assert.Equal(t, newCounterpartyID, retrievedCounterpartyID)

	assert.True(t, transaction.GetCounterpartyID().IsSome())
	assert.Equal(t, newCounterpartyID, transaction.GetCounterpartyID().Unwrap())

	// Remove counterparty
	newUpdatedAt := transaction.UpdatedAt
	time.Sleep(time.Millisecond)

	transaction.RemoveCounterparty()

	assert.False(t, transaction.HasCounterparty())
	assert.True(t, transaction.UpdatedAt.After(newUpdatedAt))

	counterpartyID, hasCounterparty = transaction.GetCounterparty()
	assert.False(t, hasCounterparty)
	assert.Equal(t, counterpartyEntity.CounterpartyID{}, counterpartyID)

	assert.True(t, transaction.GetCounterpartyID().IsNone())
}

func TestTransaction_IsDebitAndIsCredit(t *testing.T) {
	tests := []struct {
		name     string
		amount   money.Money
		isDebit  bool
		isCredit bool
	}{
		{
			name:     "positive amount is credit",
			amount:   mustMoney(t, "100.00", "USD"),
			isDebit:  false,
			isCredit: true,
		},
		{
			name:     "negative amount is debit",
			amount:   mustMoney(t, "-50.00", "USD"),
			isDebit:  true,
			isCredit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create transaction with specific amount
			transaction := &Transaction{
				Amount: tt.amount,
			}

			assert.Equal(t, tt.isDebit, transaction.IsDebit())
			assert.Equal(t, tt.isCredit, transaction.IsCredit())
		})
	}
}

func TestTransaction_GetAbsoluteAmount(t *testing.T) {
	tests := []struct {
		name           string
		amount         money.Money
		expectedAmount money.Money
	}{
		{
			name:           "positive amount",
			amount:         mustMoney(t, "100.00", "USD"),
			expectedAmount: mustMoney(t, "100.00", "USD"),
		},
		{
			name:           "negative amount",
			amount:         mustMoney(t, "-75.50", "USD"),
			expectedAmount: mustMoney(t, "75.50", "USD"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := &Transaction{
				Amount: tt.amount,
			}

			absoluteAmount := transaction.GetAbsoluteAmount()
			assert.Equal(t, tt.expectedAmount, absoluteAmount)
		})
	}
}

// Helper functions

func createTestTransaction(t *testing.T) *Transaction {
	t.Helper()

	ledgerID, err := ledgerEntity.NewLedgerID()
	require.NoError(t, err)

	accountID, err := NewAccountID()
	require.NoError(t, err)

	itemID, err := budgetEntity.NewItemID()
	require.NoError(t, err)

	amount := mustMoney(t, "100.00", "USD")

	transaction, err := NewTransaction(
		ledgerID,
		accountID,
		itemID,
		amount,
		"Test transaction",
		time.Now(),
	)
	require.NoError(t, err)

	return transaction
}