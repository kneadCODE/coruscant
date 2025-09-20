package entity

import (
	"time"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/user/entity"
)

// LedgerUser represents a user's access and role within a specific ledger
type LedgerUser struct {
	LedgerID  LedgerID
	UserID    entity.UserID
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewLedgerUser creates a new LedgerUser with specified role
func NewLedgerUser(ledgerID LedgerID, userID entity.UserID, role Role) *LedgerUser {
	now := time.Now()
	return &LedgerUser{
		LedgerID:  ledgerID,
		UserID:    userID,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ReconstructLedgerUser reconstructs a LedgerUser from stored data
func ReconstructLedgerUser(
	ledgerID LedgerID,
	userID entity.UserID,
	role Role,
	createdAt, updatedAt time.Time,
) *LedgerUser {
	return &LedgerUser{
		LedgerID:  ledgerID,
		UserID:    userID,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// UpdateRole updates the user's role in the ledger
func (lu *LedgerUser) UpdateRole(role Role) {
	lu.Role = role
	lu.UpdatedAt = time.Now()
}

// HasPermission checks if the user has a specific permission in this ledger
func (lu *LedgerUser) HasPermission(permission Permission) bool {
	return lu.Role.HasPermission(permission)
}

// CanRead checks if the user can read data in this ledger
func (lu *LedgerUser) CanRead() bool {
	return lu.HasPermission(PermissionReadOnly)
}

// CanWrite checks if the user can write data in this ledger
func (lu *LedgerUser) CanWrite() bool {
	return lu.HasPermission(PermissionEdit)
}

// IsAdmin checks if the user has admin privileges in this ledger
func (lu *LedgerUser) IsAdmin() bool {
	return lu.HasPermission(PermissionAdmin)
}
