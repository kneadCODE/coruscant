package entity

import "fmt"

// Permission represents a specific permission in the ledger
type Permission string

const (
	PermissionReadOnly Permission = "READONLY"
	PermissionEdit     Permission = "EDIT"
	PermissionAdmin    Permission = "ADMIN"
)

// NewPermission creates a new Permission from string
func NewPermission(permission string) (Permission, error) {
	switch Permission(permission) {
	case PermissionReadOnly, PermissionEdit, PermissionAdmin:
		return Permission(permission), nil
	default:
		return "", fmt.Errorf("invalid permission: %s", permission)
	}
}

// String returns the string representation of Permission
func (p Permission) String() string {
	return string(p)
}
