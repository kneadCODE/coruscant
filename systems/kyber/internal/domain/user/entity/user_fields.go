package entity

import (
	"fmt"

	"github.com/kneadCODE/coruscant/shared/golib/id"
)

// UserID represents a unique identifier for a user using UUIDv7
type UserID struct {
	id.EntityID
}

// NewUserID creates a new UserID using UUIDv7
func NewUserID() (UserID, error) {
	base, err := id.NewEntityID()
	if err != nil {
		return UserID{}, fmt.Errorf("failed to create user ID: %w", err)
	}
	return UserID{EntityID: base}, nil
}

// NewUserIDFromString creates a UserID from an existing string
func NewUserIDFromString(idStr string) (UserID, error) {
	base, err := id.NewEntityIDFromString(idStr)
	if err != nil {
		return UserID{}, fmt.Errorf("failed to create user ID: %w", err)
	}
	return UserID{EntityID: base}, nil
}

// Equals checks if two UserIDs are equal
func (u UserID) Equals(other UserID) bool {
	return u.EntityID.Equals(other.EntityID)
}
