// Package id provides entity identification using UUIDv7
package id

import (
	"fmt"
	"time"
)

// EntityID represents a base entity identifier using UUIDv7
type EntityID struct {
	value string
}

// NewEntityID creates a new EntityID using UUIDv7
func NewEntityID() (EntityID, error) {
	value, err := Generate()
	if err != nil {
		return EntityID{}, fmt.Errorf("failed to create entity ID: %w", err)
	}
	return EntityID{value: value}, nil
}

// NewEntityIDFromString creates an EntityID from an existing string
func NewEntityIDFromString(idStr string) (EntityID, error) {
	if idStr == "" {
		return EntityID{}, fmt.Errorf("entity ID cannot be empty")
	}

	validated, err := ParseID(idStr)
	if err != nil {
		return EntityID{}, fmt.Errorf("invalid entity ID format: %w", err)
	}

	return EntityID{value: validated}, nil
}

// String returns the string representation of the EntityID
func (e EntityID) String() string {
	return e.value
}

// Equals checks if two EntityIDs are equal
func (e EntityID) Equals(other EntityID) bool {
	return e.value == other.value
}

// IsValid checks if the EntityID is valid (not empty and well-formed)
func (e EntityID) IsValid() bool {
	if e.value == "" {
		return false
	}
	// Could add additional validation here (e.g., UUID format check)
	_, err := ParseID(e.value)
	return err == nil
}

// Timestamp returns the timestamp embedded in the UUIDv7
func (e EntityID) Timestamp() (time.Time, error) {
	return GetTimestamp(e.value)
}
