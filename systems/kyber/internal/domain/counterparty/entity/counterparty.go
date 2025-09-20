package entity

import (
	"fmt"
	"time"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/ledger/entity"
)

// Counterparty represents a person or organization involved in transactions
type Counterparty struct {
	ID          CounterpartyID
	LedgerID    entity.LedgerID
	Name        string
	Type        CounterpartyType
	Description string
	Status      CounterpartyStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCounterparty creates a new Counterparty
func NewCounterparty(
	ledgerID entity.LedgerID,
	name string,
	counterpartyType CounterpartyType,
	description string,
) (*Counterparty, error) {
	if !ledgerID.IsValid() {
		return nil, fmt.Errorf("ledger ID cannot be empty")
	}

	if name == "" {
		return nil, fmt.Errorf("counterparty name cannot be empty")
	}

	id, err := NewCounterpartyID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate counterparty ID: %w", err)
	}

	now := time.Now()

	return &Counterparty{
		ID:          id,
		LedgerID:    ledgerID,
		Name:        name,
		Type:        counterpartyType,
		Description: description,
		Status:      CounterpartyStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ReconstructCounterparty reconstructs a Counterparty from stored data
func ReconstructCounterparty(
	id CounterpartyID,
	ledgerID entity.LedgerID,
	name string,
	counterpartyType CounterpartyType,
	description, contactInfo string,
	status CounterpartyStatus,
	createdAt, updatedAt time.Time,
) *Counterparty {
	return &Counterparty{
		ID:          id,
		LedgerID:    ledgerID,
		Name:        name,
		Type:        counterpartyType,
		Description: description,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// UpdateInfo updates the counterparty's information
func (c *Counterparty) UpdateInfo(name, description string) error {
	if name == "" {
		return fmt.Errorf("counterparty name cannot be empty")
	}

	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now()
	return nil
}

// UpdateType updates the counterparty's type
func (c *Counterparty) UpdateType(counterpartyType CounterpartyType) {
	c.Type = counterpartyType
	c.UpdatedAt = time.Now()
}

// Activate activates the counterparty
func (c *Counterparty) Activate() {
	c.Status = CounterpartyStatusActive
	c.UpdatedAt = time.Now()
}

// Ardchive archives the counterparty
func (c *Counterparty) Ardchive() {
	c.Status = CounterpartyStatusArchived
	c.UpdatedAt = time.Now()
}
