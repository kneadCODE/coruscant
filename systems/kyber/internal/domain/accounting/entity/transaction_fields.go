package entity

import (
	"fmt"

	"github.com/kneadCODE/coruscant/shared/golib/id"
)

// TransactionID represents a unique identifier for a transaction using UUIDv7
type TransactionID struct {
	id.EntityID
}

// NewTransactionID creates a new TransactionID using UUIDv7
func NewTransactionID() (TransactionID, error) {
	base, err := id.NewEntityID()
	if err != nil {
		return TransactionID{}, fmt.Errorf("failed to create transaction ID: %w", err)
	}
	return TransactionID{EntityID: base}, nil
}

// NewTransactionIDFromString creates a TransactionID from an existing string
func NewTransactionIDFromString(idStr string) (TransactionID, error) {
	base, err := id.NewEntityIDFromString(idStr)
	if err != nil {
		return TransactionID{}, fmt.Errorf("failed to create transaction ID: %w", err)
	}
	return TransactionID{EntityID: base}, nil
}

// Equals checks if two TransactionIDs are equal
func (t TransactionID) Equals(other TransactionID) bool {
	return t.EntityID.Equals(other.EntityID)
}
