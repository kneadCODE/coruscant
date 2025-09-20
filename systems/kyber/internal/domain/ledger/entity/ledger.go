package entity

import (
	"fmt"
	"time"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/user/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

// Ledger represents the main aggregate root for a user's financial ledger
type Ledger struct {
	ID           LedgerID
	Name         string
	Description  string
	BaseCurrency money.Currency
	Status       LedgerStatus
	Users        []LedgerUser // RBAC: Users with access to this ledger (includes owner)
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewLedger creates a new Ledger with the owner as admin
func NewLedger(name, description string, baseCurrency money.Currency, adminUserID entity.UserID) (*Ledger, error) {
	if name == "" {
		return nil, fmt.Errorf("ledger name cannot be empty")
	}

	if baseCurrency == "" {
		return nil, fmt.Errorf("base currency cannot be empty")
	}

	if !adminUserID.IsValid() {
		return nil, fmt.Errorf("admin user ID is invalid")
	}

	ledgerID, err := NewLedgerID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ledger ID: %w", err)
	}

	now := time.Now()
	return &Ledger{
		ID:           ledgerID,
		Name:         name,
		Description:  description,
		BaseCurrency: baseCurrency,
		Status:       LedgerStatusActive,
		Users:        []LedgerUser{*NewLedgerUser(ledgerID, adminUserID, RoleAdmin)},
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// ReconstructLedger reconstructs a Ledger from stored data
func ReconstructLedger(
	id LedgerID,
	name, description string,
	baseCurrency money.Currency,
	status LedgerStatus,
	users []LedgerUser,
	createdAt, updatedAt time.Time,
) *Ledger {
	return &Ledger{
		ID:           id,
		Name:         name,
		Description:  description,
		BaseCurrency: baseCurrency,
		Status:       status,
		Users:        users,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

// UpdateInfo updates the ledger's basic information
func (l *Ledger) UpdateInfo(name, description string) error {
	if name == "" {
		return fmt.Errorf("ledger name cannot be empty")
	}

	l.Name = name
	l.Description = description
	l.UpdatedAt = time.Now()
	return nil
}

// UpdateUserRole updates a user's role in the ledger
func (l *Ledger) UpdateUserRole(userID entity.UserID, role Role) error {
	// TODO: Ensure there is always one admin in the ledger

	for i := range l.Users {
		if l.Users[i].UserID.Equals(userID) {
			l.Users[i].UpdateRole(role)
			l.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("user not found in ledger")
}

// GetUserAccess returns the user's access information for this ledger
func (l *Ledger) GetUserAccess(userID entity.UserID) (*LedgerUser, error) {
	for i := range l.Users {
		if l.Users[i].UserID.Equals(userID) {
			return &l.Users[i], nil
		}
	}

	return nil, fmt.Errorf("user does not have access to this ledger")
}

// HasUserAccess checks if a user has access to this ledger
func (l *Ledger) HasUserAccess(userID entity.UserID) bool {
	_, err := l.GetUserAccess(userID)
	return err == nil
}

// UserHasPermission checks if a user has a specific permission in this ledger
func (l *Ledger) UserHasPermission(userID entity.UserID, permission Permission) bool {
	userAccess, err := l.GetUserAccess(userID)
	if err != nil {
		return false
	}

	return userAccess.HasPermission(permission)
}

// GetOwner returns the admin of this ledger
func (l *Ledger) GetAdmin() *LedgerUser {
	for i := range l.Users {
		if l.Users[i].Role.HasPermission(PermissionAdmin) {
			return &l.Users[i]
		}
	}
	return nil
}

// Archive changes the ledger status to archived (read-only)
func (l *Ledger) Archive() error {
	if l.Status == LedgerStatusArchived {
		return fmt.Errorf("ledger is already archived")
	}
	l.Status = LedgerStatusArchived
	l.UpdatedAt = time.Now()
	return nil
}

// Activate changes the ledger status to active
func (l *Ledger) Activate() error {
	if l.Status == LedgerStatusActive {
		return fmt.Errorf("ledger is already active")
	}
	l.Status = LedgerStatusActive
	l.UpdatedAt = time.Now()
	return nil
}

// CanWrite checks if the ledger allows write operations
func (l *Ledger) CanWrite() bool {
	return l.Status.CanWrite()
}

// CanRead checks if the ledger allows read operations
func (l *Ledger) CanRead() bool {
	return l.Status.CanRead()
}

// UserCount returns the number of users with access to this ledger
func (l *Ledger) UserCount() int {
	return len(l.Users)
}
