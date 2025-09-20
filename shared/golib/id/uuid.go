// Package id provides UUID generation and validation utilities
package id

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Generate generates a new UUIDv7
func Generate() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUIDv7: %w", err)
	}
	return id.String(), nil
}

// ParseID validates a UUID string
func ParseID(id string) (string, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("invalid UUID: %w", err)
	}
	return parsed.String(), nil
}

// GetTimestamp extracts the timestamp from a UUIDv7
func GetTimestamp(id string) (time.Time, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid UUID: %w", err)
	}

	// UUIDv7 timestamp is in the first 48 bits (6 bytes)
	timestampMS := int64(parsed[0])<<40 | int64(parsed[1])<<32 | int64(parsed[2])<<24 |
		int64(parsed[3])<<16 | int64(parsed[4])<<8 | int64(parsed[5])

	return time.Unix(timestampMS/1000, (timestampMS%1000)*1000000), nil
}
