package entity

import (
	"fmt"

	"github.com/kneadCODE/coruscant/shared/golib/id"
)

// LedgerID represents a unique identifier for a ledger using UUIDv7
type LedgerID struct {
	id.EntityID
}

// NewLedgerID creates a new LedgerID using UUIDv7
func NewLedgerID() (LedgerID, error) {
	base, err := id.NewEntityID()
	if err != nil {
		return LedgerID{}, fmt.Errorf("failed to create ledger ID: %w", err)
	}
	return LedgerID{EntityID: base}, nil
}

// NewLedgerIDFromString creates a LedgerID from an existing string
func NewLedgerIDFromString(idStr string) (LedgerID, error) {
	base, err := id.NewEntityIDFromString(idStr)
	if err != nil {
		return LedgerID{}, fmt.Errorf("failed to create ledger ID: %w", err)
	}
	return LedgerID{EntityID: base}, nil
}

// Equals checks if two LedgerIDs are equal
func (l LedgerID) Equals(other LedgerID) bool {
	return l.EntityID.Equals(other.EntityID)
}

// LedgerStatus represents the current status of a ledger
type LedgerStatus string

const (
	LedgerStatusActive   LedgerStatus = "ACTIVE"   // Active ledger, normal operations
	LedgerStatusArchived LedgerStatus = "ARCHIVED" // Archived ledger, read-only access
)

// IsActive checks if the ledger status allows normal operations
func (s LedgerStatus) IsActive() bool {
	return s == LedgerStatusActive
}

// IsArchived checks if the ledger is archived (read-only)
func (s LedgerStatus) IsArchived() bool {
	return s == LedgerStatusArchived
}

// CanWrite checks if the ledger allows write operations
func (s LedgerStatus) CanWrite() bool {
	return s.IsActive()
}

// CanRead checks if the ledger allows read operations
func (s LedgerStatus) CanRead() bool {
	return s == LedgerStatusActive || s == LedgerStatusArchived
}
