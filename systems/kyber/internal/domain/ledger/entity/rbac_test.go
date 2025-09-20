package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermission_NewPermission(t *testing.T) {
	tests := []struct {
		name        string
		permission  string
		expected    Permission
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid readonly permission",
			permission: "READONLY",
			expected:   PermissionReadOnly,
			wantErr:    false,
		},
		{
			name:       "valid edit permission",
			permission: "EDIT",
			expected:   PermissionEdit,
			wantErr:    false,
		},
		{
			name:       "valid admin permission",
			permission: "ADMIN",
			expected:   PermissionAdmin,
			wantErr:    false,
		},
		{
			name:        "invalid permission",
			permission:  "INVALID",
			expected:    "",
			wantErr:     true,
			errContains: "invalid permission",
		},
		{
			name:        "empty permission",
			permission:  "",
			expected:    "",
			wantErr:     true,
			errContains: "invalid permission",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission, err := NewPermission(tt.permission)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, Permission(""), permission)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, permission)
			}
		})
	}
}

func TestPermission_String(t *testing.T) {
	tests := []struct {
		name       string
		permission Permission
		expected   string
	}{
		{
			name:       "readonly permission string",
			permission: PermissionReadOnly,
			expected:   "READONLY",
		},
		{
			name:       "edit permission string",
			permission: PermissionEdit,
			expected:   "EDIT",
		},
		{
			name:       "admin permission string",
			permission: PermissionAdmin,
			expected:   "ADMIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.permission.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewRole(t *testing.T) {
	tests := []struct {
		name        string
		roleName    string
		permissions []Permission
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid custom role",
			roleName:    "CUSTOM_ROLE",
			permissions: []Permission{PermissionReadOnly, PermissionEdit},
			wantErr:     false,
		},
		{
			name:        "empty role name",
			roleName:    "",
			permissions: []Permission{PermissionReadOnly},
			wantErr:     true,
			errContains: "role name cannot be empty",
		},
		{
			name:        "empty permissions",
			roleName:    "EMPTY_ROLE",
			permissions: []Permission{},
			wantErr:     true,
			errContains: "role must have at least one permission",
		},
		{
			name:        "nil permissions",
			roleName:    "NIL_ROLE",
			permissions: nil,
			wantErr:     true,
			errContains: "role must have at least one permission",
		},
		{
			name:        "single permission role",
			roleName:    "READONLY_ROLE",
			permissions: []Permission{PermissionReadOnly},
			wantErr:     false,
		},
		{
			name:        "all permissions role",
			roleName:    "FULL_ACCESS",
			permissions: []Permission{PermissionReadOnly, PermissionEdit, PermissionAdmin},
			wantErr:     false,
		},
		{
			name:        "duplicate permissions (should work)",
			roleName:    "DUPLICATE_PERMS",
			permissions: []Permission{PermissionReadOnly, PermissionReadOnly, PermissionEdit},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := NewRole(tt.roleName, tt.permissions)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Equal(t, Role{}, role)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.roleName, role.Name)
				assert.NotNil(t, role.Permissions)

				// Verify all permissions are present
				for _, perm := range tt.permissions {
					assert.True(t, role.HasPermission(perm))
				}

				// For duplicate permissions test, verify map behavior
				if tt.name == "duplicate permissions (should work)" {
					assert.Len(t, role.Permissions, 2) // Should have 2 unique permissions
				}
			}
		})
	}
}

func TestRole_HasPermission(t *testing.T) {
	// Test predefined roles
	tests := []struct {
		name       string
		role       Role
		permission Permission
		expected   bool
	}{
		// RoleViewer tests
		{
			name:       "viewer has readonly permission",
			role:       RoleViewer,
			permission: PermissionReadOnly,
			expected:   true,
		},
		{
			name:       "viewer does not have edit permission",
			role:       RoleViewer,
			permission: PermissionEdit,
			expected:   false,
		},
		{
			name:       "viewer does not have admin permission",
			role:       RoleViewer,
			permission: PermissionAdmin,
			expected:   false,
		},

		// RoleEditor tests
		{
			name:       "editor has readonly permission",
			role:       RoleEditor,
			permission: PermissionReadOnly,
			expected:   true,
		},
		{
			name:       "editor has edit permission",
			role:       RoleEditor,
			permission: PermissionEdit,
			expected:   true,
		},
		{
			name:       "editor does not have admin permission",
			role:       RoleEditor,
			permission: PermissionAdmin,
			expected:   false,
		},

		// RoleAdmin tests
		{
			name:       "admin has readonly permission",
			role:       RoleAdmin,
			permission: PermissionReadOnly,
			expected:   true,
		},
		{
			name:       "admin has edit permission",
			role:       RoleAdmin,
			permission: PermissionEdit,
			expected:   true,
		},
		{
			name:       "admin has admin permission",
			role:       RoleAdmin,
			permission: PermissionAdmin,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test custom role
	t.Run("custom role permissions", func(t *testing.T) {
		customRole, err := NewRole("CUSTOM", []Permission{PermissionReadOnly, PermissionAdmin})
		require.NoError(t, err)

		assert.True(t, customRole.HasPermission(PermissionReadOnly))
		assert.False(t, customRole.HasPermission(PermissionEdit))
		assert.True(t, customRole.HasPermission(PermissionAdmin))
	})

	// Test empty role
	t.Run("empty role has no permissions", func(t *testing.T) {
		emptyRole := Role{
			Name:        "EMPTY",
			Permissions: make(map[Permission]bool),
		}

		assert.False(t, emptyRole.HasPermission(PermissionReadOnly))
		assert.False(t, emptyRole.HasPermission(PermissionEdit))
		assert.False(t, emptyRole.HasPermission(PermissionAdmin))
	})
}

func TestRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected string
	}{
		{
			name:     "viewer role string",
			role:     RoleViewer,
			expected: "VIEWER",
		},
		{
			name:     "editor role string",
			role:     RoleEditor,
			expected: "EDITOR",
		},
		{
			name:     "admin role string",
			role:     RoleAdmin,
			expected: "ADMIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.String()
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test custom role
	t.Run("custom role string", func(t *testing.T) {
		customRole, err := NewRole("CUSTOM_ROLE", []Permission{PermissionReadOnly})
		require.NoError(t, err)

		result := customRole.String()
		assert.Equal(t, "CUSTOM_ROLE", result)
	})
}

func TestPredefinedRoles(t *testing.T) {
	t.Run("verify predefined roles structure", func(t *testing.T) {
		// Test RoleViewer
		assert.Equal(t, "VIEWER", RoleViewer.Name)
		assert.Len(t, RoleViewer.Permissions, 1)
		assert.True(t, RoleViewer.Permissions[PermissionReadOnly])

		// Test RoleEditor
		assert.Equal(t, "EDITOR", RoleEditor.Name)
		assert.Len(t, RoleEditor.Permissions, 2)
		assert.True(t, RoleEditor.Permissions[PermissionReadOnly])
		assert.True(t, RoleEditor.Permissions[PermissionEdit])

		// Test RoleAdmin
		assert.Equal(t, "ADMIN", RoleAdmin.Name)
		assert.Len(t, RoleAdmin.Permissions, 3)
		assert.True(t, RoleAdmin.Permissions[PermissionReadOnly])
		assert.True(t, RoleAdmin.Permissions[PermissionEdit])
		assert.True(t, RoleAdmin.Permissions[PermissionAdmin])
	})

	t.Run("role hierarchy verification", func(t *testing.T) {
		// Admin should have all permissions that Editor has
		for perm := range RoleEditor.Permissions {
			assert.True(t, RoleAdmin.HasPermission(perm), "Admin should have permission %s", perm)
		}

		// Editor should have all permissions that Viewer has
		for perm := range RoleViewer.Permissions {
			assert.True(t, RoleEditor.HasPermission(perm), "Editor should have permission %s", perm)
		}
	})
}