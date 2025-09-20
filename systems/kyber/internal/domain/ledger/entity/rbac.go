package entity

import "fmt"

// Role represents a collection of permissions
type Role struct {
	Name        string
	Permissions map[Permission]bool
}

// Predefined roles
var (
	RoleViewer = Role{
		Name: "VIEWER",
		Permissions: map[Permission]bool{
			PermissionReadOnly: true,
		},
	}

	RoleEditor = Role{
		Name: "EDITOR",
		Permissions: map[Permission]bool{
			PermissionReadOnly: true,
			PermissionEdit:     true,
		},
	}

	RoleAdmin = Role{
		Name: "ADMIN",
		Permissions: map[Permission]bool{
			PermissionReadOnly: true,
			PermissionEdit:     true,
			PermissionAdmin:    true,
		},
	}
)

// NewRole creates a new custom role
func NewRole(name string, permissions []Permission) (Role, error) {
	if name == "" {
		return Role{}, fmt.Errorf("role name cannot be empty")
	}

	if len(permissions) == 0 {
		return Role{}, fmt.Errorf("role must have at least one permission")
	}

	p := make(map[Permission]bool, len(permissions))

	for idx := range permissions {
		p[permissions[idx]] = true
	}

	return Role{
		Name:        name,
		Permissions: p,
	}, nil
}

// HasPermission checks if the role has a specific permission
func (r Role) HasPermission(permission Permission) bool {
	if _, exists := r.Permissions[permission]; exists {
		return true
	}
	return false
}

// String returns the string representation of the role
func (r Role) String() string {
	return r.Name
}
