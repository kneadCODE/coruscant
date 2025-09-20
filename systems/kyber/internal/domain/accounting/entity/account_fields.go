package entity

import (
	"fmt"

	"github.com/kneadCODE/coruscant/shared/golib/id"
)

// AccountID represents a unique identifier for an account using UUIDv7
type AccountID struct {
	id.EntityID
}

// NewAccountID creates a new AccountID using UUIDv7
func NewAccountID() (AccountID, error) {
	base, err := id.NewEntityID()
	if err != nil {
		return AccountID{}, fmt.Errorf("failed to create account ID: %w", err)
	}
	return AccountID{EntityID: base}, nil
}

// NewAccountIDFromString creates an AccountID from an existing string
func NewAccountIDFromString(idStr string) (AccountID, error) {
	base, err := id.NewEntityIDFromString(idStr)
	if err != nil {
		return AccountID{}, fmt.Errorf("failed to create account ID: %w", err)
	}
	return AccountID{EntityID: base}, nil
}

// Equals checks if two AccountIDs are equal
func (a AccountID) Equals(other AccountID) bool {
	return a.EntityID.Equals(other.EntityID)
}

// AccountType represents the type of an account
type AccountType string

const (
	AccountTypeChecking AccountType = "CHECKING"
	AccountTypeSavings  AccountType = "SAVINGS"
	AccountTypeCash     AccountType = "CASH"

	AccountTypeInvestment    AccountType = "INVESTMENT"
	AccountTypeDigitalWallet AccountType = "DIGITAL_WALLET"

	AccountTypeHolding AccountType = "HOLDING"

	AccountTypeCreditCard  AccountType = "CREDIT_CARD"
	AccountTypeInstallment AccountType = "INSTALLMENT"
	AccountTypeLoan        AccountType = "LOAN"
	AccountTypeMortgage    AccountType = "MORTGAGE"
)

// NewAccountType creates a new AccountType from string
func NewAccountType(accountType string) (AccountType, error) {
	switch AccountType(accountType) {
	case AccountTypeSavings,
		AccountTypeChecking,
		AccountTypeCash,
		AccountTypeHolding,
		AccountTypeInvestment,
		AccountTypeDigitalWallet,
		AccountTypeCreditCard,
		AccountTypeInstallment,
		AccountTypeLoan,
		AccountTypeMortgage:
		return AccountType(accountType), nil
	default:
		return "", fmt.Errorf("invalid account type: %s", accountType)
	}
}

// String returns the string representation of AccountType
func (a AccountType) String() string {
	return string(a)
}

// IsAsset checks if the account type is an asset
func (a AccountType) IsAsset() bool {
	switch a {
	case AccountTypeSavings,
		AccountTypeChecking,
		AccountTypeCash,
		AccountTypeHolding,
		AccountTypeInvestment,
		AccountTypeDigitalWallet:
		return true
	default:
		return false
	}
}

// IsLiability checks if the account type is a liability
func (a AccountType) IsLiability() bool {
	switch a {
	case AccountTypeCreditCard,
		AccountTypeInstallment,
		AccountTypeLoan,
		AccountTypeMortgage:
		return true
	default:
		return false
	}
}

// func (a AccountType) IsEquity() bool {
// }

// AccountStatus represents the current status of a account
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "ACTIVE"   // Active account, normal operations
	AccountStatusArchived AccountStatus = "ARCHIVED" // Archived account, read-only access
)

// IsActive checks if the account status allows normal operations
func (s AccountStatus) IsActive() bool {
	return s == AccountStatusActive
}

// IsArchived checks if the account is archived (read-only)
func (s AccountStatus) IsArchived() bool {
	return s == AccountStatusArchived
}

// AccountCategory represents the accounting category of an account
type AccountCategory string

const (
	AccountCategoryAsset     AccountCategory = "ASSET"     // Assets (positive balance = more valuable)
	AccountCategoryLiability AccountCategory = "LIABILITY" // Liabilities (positive balance = more owed)
	AccountCategoryEquity    AccountCategory = "EQUITY"    // Equity (ownership/net worth)
)

// NewAccountCategory creates a new AccountCategory from string
func NewAccountCategory(category string) (AccountCategory, error) {
	switch AccountCategory(category) {
	case AccountCategoryAsset, AccountCategoryLiability, AccountCategoryEquity:
		return AccountCategory(category), nil
	default:
		return "", fmt.Errorf("invalid account category: %s", category)
	}
}

// String returns the string representation of AccountCategory
func (a AccountCategory) String() string {
	return string(a)
}

// IsAsset checks if the category is an asset
func (a AccountCategory) IsAsset() bool {
	return a == AccountCategoryAsset
}

// IsLiability checks if the category is a liability
func (a AccountCategory) IsLiability() bool {
	return a == AccountCategoryLiability
}

// IsEquity checks if the category is equity
func (a AccountCategory) IsEquity() bool {
	return a == AccountCategoryEquity
}

// AllAccountCategories returns all valid account categories
func AllAccountCategories() []AccountCategory {
	return []AccountCategory{
		AccountCategoryAsset,
		AccountCategoryLiability,
		AccountCategoryEquity,
	}
}
