package entity

import (
	"fmt"

	"github.com/kneadCODE/coruscant/shared/golib/id"
)

// CounterpartyID represents a unique identifier for a counterparty using UUIDv7
type CounterpartyID struct {
	id.EntityID
}

// NewCounterpartyID creates a new CounterpartyID using UUIDv7
func NewCounterpartyID() (CounterpartyID, error) {
	base, err := id.NewEntityID()
	if err != nil {
		return CounterpartyID{}, fmt.Errorf("failed to create counterparty ID: %w", err)
	}
	return CounterpartyID{EntityID: base}, nil
}

// NewCounterpartyIDFromString creates a CounterpartyID from an existing string
func NewCounterpartyIDFromString(idStr string) (CounterpartyID, error) {
	base, err := id.NewEntityIDFromString(idStr)
	if err != nil {
		return CounterpartyID{}, fmt.Errorf("failed to create counterparty ID: %w", err)
	}
	return CounterpartyID{EntityID: base}, nil
}

// Equals checks if two CounterpartyIDs are equal
func (c CounterpartyID) Equals(other CounterpartyID) bool {
	return c.EntityID.Equals(other.EntityID)
}

// CounterpartyType represents the type of counterparty
type CounterpartyType string

// Counterparty type constants define different types of transaction counterparties
const (
	CounterpartyTypeIndividual    CounterpartyType = "INDIVIDUAL"     // Individual person
	CounterpartyTypeOrganization  CounterpartyType = "ORGANIZATION"   // General organization
	CounterpartyTypeBusiness      CounterpartyType = "BUSINESS"       // Business entity
	CounterpartyTypeGovernment    CounterpartyType = "GOVERNMENT"     // Government entity
	CounterpartyTypeUtility       CounterpartyType = "UTILITY"        // Utility company
	CounterpartyTypeOnlineService CounterpartyType = "ONLINE_SERVICE" // Online service provider
	CounterpartyTypeRetailer      CounterpartyType = "RETAILER"       // Retail business
	CounterpartyTypeInvestment    CounterpartyType = "INVESTMENT"     // Investment firm
	CounterpartyTypeInsurance     CounterpartyType = "INSURANCE"      // Insurance company
	CounterpartyTypeEmployer      CounterpartyType = "EMPLOYER"       // Employer organization
	CounterpartyTypeBank          CounterpartyType = "BANK"           // Bank or financial institution
)

// NewCounterpartyType creates a new CounterpartyType from string
func NewCounterpartyType(counterpartyType string) (CounterpartyType, error) {
	switch CounterpartyType(counterpartyType) {
	case CounterpartyTypeIndividual, CounterpartyTypeOrganization, CounterpartyTypeBusiness,
		CounterpartyTypeGovernment, CounterpartyTypeUtility, CounterpartyTypeOnlineService,
		CounterpartyTypeRetailer, CounterpartyTypeInvestment, CounterpartyTypeInsurance,
		CounterpartyTypeEmployer, CounterpartyTypeBank:
		return CounterpartyType(counterpartyType), nil
	default:
		return "", fmt.Errorf("invalid counterparty: %s", counterpartyType)
	}
}

// String returns the string representation of CounterpartyType
func (c CounterpartyType) String() string {
	return string(c)
}

// CounterpartyStatus represents the current status of a counterparty
type CounterpartyStatus string

// Counterparty status constants define the operational state of counterparties
const (
	CounterpartyStatusActive   CounterpartyStatus = "ACTIVE"   // Active counterparty, normal operations
	CounterpartyStatusArchived CounterpartyStatus = "ARCHIVED" // Archived counterparty, read-only access
)

// IsActive checks if the counterparty status allows normal operations
func (s CounterpartyStatus) IsActive() bool {
	return s == CounterpartyStatusActive
}

// IsArchived checks if the counterparty is archived (read-only)
func (s CounterpartyStatus) IsArchived() bool {
	return s == CounterpartyStatusArchived
}

// AllCounterpartyTypes returns all valid counterparty types
func AllCounterpartyTypes() []CounterpartyType {
	return []CounterpartyType{
		CounterpartyTypeIndividual,
		CounterpartyTypeBusiness,
		CounterpartyTypeOrganization,
		CounterpartyTypeGovernment,
		CounterpartyTypeUtility,
		CounterpartyTypeOnlineService,
		CounterpartyTypeRetailer,
		CounterpartyTypeInvestment,
		CounterpartyTypeInsurance,
		CounterpartyTypeEmployer,
		CounterpartyTypeBank,
	}
}
