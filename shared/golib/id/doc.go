// Package id provides a unified entity identification system using UUIDv7.
//
// This package consolidates ID generation and validation for all domain entities.
// It uses UUIDv7 which provides time-ordered, sortable identifiers with embedded timestamps.
//
// EntityID serves as a base type that other domain-specific ID types can embed
// for type safety while sharing common functionality.
//
// Basic Usage:
//
//	// Generate a new ID
//	id, err := id.NewEntityID()
//	if err != nil {
//		return err
//	}
//
//	// Create from existing string
//	id, err := id.NewEntityIDFromString("01234567-89ab-cdef-0123-456789abcdef")
//	if err != nil {
//		return err
//	}
//
//	// Validation
//	if id.IsValid() {
//		// ID is well-formed
//	}
//
//	// String representation
//	idStr := id.String()
//
// UUIDv7 Features:
//
//	UUIDv7 provides several advantages over other UUID versions:
//	- Time-ordered: IDs sort chronologically
//	- Embedded timestamp: Extract creation time
//	- Database-friendly: Better indexing performance
//	- Globally unique: No coordination required
//
// Timestamp Extraction:
//
//	// Get the embedded timestamp
//	timestamp, err := id.Timestamp()
//	if err != nil {
//		return err
//	}
//	fmt.Printf("ID created at: %s", timestamp)
//
// Domain Integration:
//
//	Domain-specific ID types embed EntityID for type safety:
//
//	type UserID struct {
//		id.EntityID
//	}
//
//	func NewUserID() (UserID, error) {
//		base, err := id.NewEntityID()
//		if err != nil {
//			return UserID{}, err
//		}
//		return UserID{EntityID: base}, nil
//	}
//
// This approach provides strong typing while eliminating code duplication
// across different entity types.
package id
