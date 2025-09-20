package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

func TestNewAccount(t *testing.T) {
	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	tests := []struct {
		name        string
		ledgerID    entity.LedgerID
		accountName string
		description string
		accountType AccountType
		currency    money.Currency
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid account creation",
			ledgerID:    ledgerID,
			accountName: "Checking Account",
			description: "Primary checking account",
			accountType: AccountTypeChecking,
			currency:    "USD",
			wantErr:     false,
		},
		{
			name:        "invalid ledger ID",
			ledgerID:    entity.LedgerID{},
			accountName: "Test Account",
			description: "Test description",
			accountType: AccountTypeChecking,
			currency:    "USD",
			wantErr:     true,
			errContains: "ledger ID is invalid",
		},
		{
			name:        "empty account name",
			ledgerID:    ledgerID,
			accountName: "",
			description: "Test description",
			accountType: AccountTypeChecking,
			currency:    "USD",
			wantErr:     true,
			errContains: "account name cannot be empty",
		},
		{
			name:        "empty currency",
			ledgerID:    ledgerID,
			accountName: "Test Account",
			description: "Test description",
			accountType: AccountTypeChecking,
			currency:    "",
			wantErr:     true,
			errContains: "account currency cannot be empty",
		},
		{
			name:        "valid credit card account",
			ledgerID:    ledgerID,
			accountName: "Credit Card",
			description: "Primary credit card",
			accountType: AccountTypeCreditCard,
			currency:    "USD",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccount(tt.ledgerID, tt.accountName, tt.description, tt.accountType, tt.currency)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, account)
			} else {
				require.NoError(t, err)
				require.NotNil(t, account)

				assert.True(t, account.ID.IsValid())
				assert.Equal(t, tt.ledgerID, account.LedgerID)
				assert.Equal(t, tt.accountName, account.Name)
				assert.Equal(t, tt.description, account.Description)
				assert.Equal(t, tt.accountType, account.Type)
				assert.Equal(t, tt.currency, account.Currency)
				assert.Equal(t, AccountStatusActive, account.Status)
				assert.True(t, account.Balance.IsZero())
				assert.Equal(t, tt.currency, account.Balance.Currency)
				assert.False(t, account.CreatedAt.IsZero())
				assert.False(t, account.UpdatedAt.IsZero())
			}
		})
	}
}

func TestReconstructAccount(t *testing.T) {
	accountID, err := NewAccountID()
	require.NoError(t, err)

	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	balance, err := money.NewMoney("100.50", "USD")
	require.NoError(t, err)

	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	account := ReconstructAccount(
		accountID,
		ledgerID,
		"Test Account",
		"Test description",
		AccountTypeChecking,
		"USD",
		balance,
		AccountStatusArchived,
		createdAt,
		updatedAt,
	)

	assert.Equal(t, accountID, account.ID)
	assert.Equal(t, ledgerID, account.LedgerID)
	assert.Equal(t, "Test Account", account.Name)
	assert.Equal(t, "Test description", account.Description)
	assert.Equal(t, AccountTypeChecking, account.Type)
	assert.Equal(t, money.Currency("USD"), account.Currency)
	assert.Equal(t, balance, account.Balance)
	assert.Equal(t, AccountStatusArchived, account.Status)
	assert.Equal(t, createdAt, account.CreatedAt)
	assert.Equal(t, updatedAt, account.UpdatedAt)
}

func TestAccount_UpdateInfo(t *testing.T) {
	account := createTestAccount(t)
	originalUpdatedAt := account.UpdatedAt

	time.Sleep(time.Millisecond) // Ensure time difference

	tests := []struct {
		name        string
		newName     string
		description string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid update",
			newName:     "Updated Account Name",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "empty name",
			newName:     "",
			description: "Some description",
			wantErr:     true,
			errContains: "account name cannot be empty",
		},
		{
			name:        "empty description is allowed",
			newName:     "Account Name",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := account.UpdateInfo(tt.newName, tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newName, account.Name)
				assert.Equal(t, tt.description, account.Description)
				assert.True(t, account.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestAccount_Credit(t *testing.T) {
	account := createTestAccount(t)
	originalBalance := account.Balance

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid credit",
			amount:  mustMoney(t, "50.00", "USD"),
			wantErr: false,
		},
		{
			name:        "currency mismatch",
			amount:      mustMoney(t, "50.00", "EUR"),
			wantErr:     true,
			errContains: "currency mismatch",
		},
		{
			name:        "negative amount",
			amount:      mustMoney(t, "-50.00", "USD"),
			wantErr:     true,
			errContains: "credit amount cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeBalance := account.Balance
			beforeUpdatedAt := account.UpdatedAt

			time.Sleep(time.Millisecond)
			err := account.Credit(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, beforeBalance, account.Balance)
			} else {
				require.NoError(t, err)
				expectedBalance, _ := beforeBalance.Add(tt.amount)
				assert.Equal(t, expectedBalance, account.Balance)
				assert.True(t, account.UpdatedAt.After(beforeUpdatedAt))
			}
		})
	}

	// Verify the account was credited successfully
	if !originalBalance.IsZero() {
		expectedFinalBalance, _ := originalBalance.Add(mustMoney(t, "50.00", "USD"))
		assert.Equal(t, expectedFinalBalance, account.Balance)
	}
}

func TestAccount_Debit(t *testing.T) {
	account := createTestAccount(t)
	// Add some balance first
	err := account.Credit(mustMoney(t, "100.00", "USD"))
	require.NoError(t, err)

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid debit",
			amount:  mustMoney(t, "30.00", "USD"),
			wantErr: false,
		},
		{
			name:        "currency mismatch",
			amount:      mustMoney(t, "50.00", "EUR"),
			wantErr:     true,
			errContains: "currency mismatch",
		},
		{
			name:        "negative amount",
			amount:      mustMoney(t, "-25.00", "USD"),
			wantErr:     true,
			errContains: "debit amount cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeBalance := account.Balance
			beforeUpdatedAt := account.UpdatedAt

			time.Sleep(time.Millisecond)
			err := account.Debit(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, beforeBalance, account.Balance)
			} else {
				require.NoError(t, err)
				expectedBalance, _ := beforeBalance.Subtract(tt.amount)
				assert.Equal(t, expectedBalance, account.Balance)
				assert.True(t, account.UpdatedAt.After(beforeUpdatedAt))
			}
		})
	}
}

func TestAccount_ActivateAndArchive(t *testing.T) {
	account := createTestAccount(t)
	account.Status = AccountStatusArchived
	originalUpdatedAt := account.UpdatedAt

	time.Sleep(time.Millisecond)

	// Test Activate
	account.Activate()
	assert.Equal(t, AccountStatusActive, account.Status)
	assert.True(t, account.UpdatedAt.After(originalUpdatedAt))

	// Test Archive
	newUpdatedAt := account.UpdatedAt
	time.Sleep(time.Millisecond)
	account.Archive()
	assert.Equal(t, AccountStatusArchived, account.Status)
	assert.True(t, account.UpdatedAt.After(newUpdatedAt))
}

func TestAccount_CanDebit(t *testing.T) {
	// Test asset account (checking)
	assetAccount := createTestAccount(t)
	err := assetAccount.Credit(mustMoney(t, "100.00", "USD"))
	require.NoError(t, err)

	// Test liability account (credit card)
	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)
	liabilityAccount, err := NewAccount(ledgerID, "Credit Card", "Test credit card", AccountTypeCreditCard, "USD")
	require.NoError(t, err)

	tests := []struct {
		name     string
		account  *Account
		amount   money.Money
		expected bool
	}{
		{
			name:     "asset account with sufficient balance",
			account:  assetAccount,
			amount:   mustMoney(t, "50.00", "USD"),
			expected: true,
		},
		{
			name:     "asset account with insufficient balance",
			account:  assetAccount,
			amount:   mustMoney(t, "150.00", "USD"),
			expected: false,
		},
		{
			name:     "liability account allows any positive amount",
			account:  liabilityAccount,
			amount:   mustMoney(t, "1000.00", "USD"),
			expected: true,
		},
		{
			name:     "currency mismatch",
			account:  assetAccount,
			amount:   mustMoney(t, "50.00", "EUR"),
			expected: false,
		},
		{
			name:     "negative amount",
			account:  assetAccount,
			amount:   mustMoney(t, "-50.00", "USD"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.account.CanDebit(tt.amount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAccount_HasSufficientBalance(t *testing.T) {
	account := createTestAccount(t)
	err := account.Credit(mustMoney(t, "100.00", "USD"))
	require.NoError(t, err)

	// HasSufficientBalance should be an alias for CanDebit
	amount := mustMoney(t, "50.00", "USD")
	assert.Equal(t, account.CanDebit(amount), account.HasSufficientBalance(amount))
}

func TestAccount_GetBalanceFloat64(t *testing.T) {
	account := createTestAccount(t)
	err := account.Credit(mustMoney(t, "123.45", "USD"))
	require.NoError(t, err)

	balance := account.GetBalanceFloat64()
	assert.Equal(t, 123.45, balance)
}

func TestAccount_DebitBalance(t *testing.T) {
	account := createTestAccount(t)
	err := account.Credit(mustMoney(t, "100.00", "USD"))
	require.NoError(t, err)

	tests := []struct {
		name        string
		amount      money.Money
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid debit with sufficient balance",
			amount:  mustMoney(t, "30.00", "USD"),
			wantErr: false,
		},
		{
			name:        "insufficient balance",
			amount:      mustMoney(t, "1000.00", "USD"),
			wantErr:     true,
			errContains: "insufficient balance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeBalance := account.Balance
			err := account.DebitBalance(tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, beforeBalance, account.Balance)
			} else {
				require.NoError(t, err)
				expectedBalance, _ := beforeBalance.Subtract(tt.amount)
				assert.Equal(t, expectedBalance, account.Balance)
			}
		})
	}
}

func TestAccount_CreditBalance(t *testing.T) {
	account := createTestAccount(t)
	beforeBalance := account.Balance

	amount := mustMoney(t, "75.50", "USD")
	err := account.CreditBalance(amount)

	require.NoError(t, err)
	expectedBalance, _ := beforeBalance.Add(amount)
	assert.Equal(t, expectedBalance, account.Balance)
}

// Helper functions

func createTestAccount(t *testing.T) *Account {
	t.Helper()

	ledgerID, err := entity.NewLedgerID()
	require.NoError(t, err)

	account, err := NewAccount(ledgerID, "Test Account", "Test description", AccountTypeChecking, "USD")
	require.NoError(t, err)

	return account
}

func mustMoney(t *testing.T, amount, currency string) money.Money {
	t.Helper()

	m, err := money.NewMoney(amount, money.Currency(currency))
	require.NoError(t, err)
	return m
}