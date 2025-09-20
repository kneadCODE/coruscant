package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/user/entity"
)

func TestNewLedgerUser(t *testing.T) {
	ledgerID, err := NewLedgerID()
	require.NoError(t, err)

	userID, err := entity.NewUserID()
	require.NoError(t, err)

	ledgerUser := NewLedgerUser(ledgerID, userID, RoleEditor)

	assert.Equal(t, ledgerID, ledgerUser.LedgerID)
	assert.Equal(t, userID, ledgerUser.UserID)
	assert.Equal(t, RoleEditor, ledgerUser.Role)
	assert.False(t, ledgerUser.CreatedAt.IsZero())
	assert.False(t, ledgerUser.UpdatedAt.IsZero())
	assert.Equal(t, ledgerUser.CreatedAt, ledgerUser.UpdatedAt)
}

func TestReconstructLedgerUser(t *testing.T) {
	ledgerID, err := NewLedgerID()
	require.NoError(t, err)

	userID, err := entity.NewUserID()
	require.NoError(t, err)

	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	ledgerUser := ReconstructLedgerUser(
		ledgerID,
		userID,
		RoleViewer,
		createdAt,
		updatedAt,
	)

	assert.Equal(t, ledgerID, ledgerUser.LedgerID)
	assert.Equal(t, userID, ledgerUser.UserID)
	assert.Equal(t, RoleViewer, ledgerUser.Role)
	assert.Equal(t, createdAt, ledgerUser.CreatedAt)
	assert.Equal(t, updatedAt, ledgerUser.UpdatedAt)
}

func TestLedgerUser_UpdateRole(t *testing.T) {
	ledgerUser := createTestLedgerUser(t)
	originalUpdatedAt := ledgerUser.UpdatedAt

	time.Sleep(time.Millisecond)

	ledgerUser.UpdateRole(RoleAdmin)

	assert.Equal(t, RoleAdmin, ledgerUser.Role)
	assert.True(t, ledgerUser.UpdatedAt.After(originalUpdatedAt))
}

func TestLedgerUser_HasPermission(t *testing.T) {
	tests := []struct {
		name       string
		role       Role
		permission Permission
		expected   bool
	}{
		{
			name:       "admin has admin permission",
			role:       RoleAdmin,
			permission: PermissionAdmin,
			expected:   true,
		},
		{
			name:       "admin has edit permission",
			role:       RoleAdmin,
			permission: PermissionEdit,
			expected:   true,
		},
		{
			name:       "admin has read permission",
			role:       RoleAdmin,
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
			name:       "editor has read permission",
			role:       RoleEditor,
			permission: PermissionReadOnly,
			expected:   true,
		},
		{
			name:       "editor does not have admin permission",
			role:       RoleEditor,
			permission: PermissionAdmin,
			expected:   false,
		},
		{
			name:       "viewer has read permission",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledgerUser := createTestLedgerUser(t)
			ledgerUser.Role = tt.role

			result := ledgerUser.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLedgerUser_CanRead(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{
			name:     "admin can read",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "editor can read",
			role:     RoleEditor,
			expected: true,
		},
		{
			name:     "viewer can read",
			role:     RoleViewer,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledgerUser := createTestLedgerUser(t)
			ledgerUser.Role = tt.role

			result := ledgerUser.CanRead()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLedgerUser_CanWrite(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{
			name:     "admin can write",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "editor can write",
			role:     RoleEditor,
			expected: true,
		},
		{
			name:     "viewer cannot write",
			role:     RoleViewer,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledgerUser := createTestLedgerUser(t)
			ledgerUser.Role = tt.role

			result := ledgerUser.CanWrite()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLedgerUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{
			name:     "admin is admin",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "editor is not admin",
			role:     RoleEditor,
			expected: false,
		},
		{
			name:     "viewer is not admin",
			role:     RoleViewer,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledgerUser := createTestLedgerUser(t)
			ledgerUser.Role = tt.role

			result := ledgerUser.IsAdmin()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions

func createTestLedgerUser(t *testing.T) *LedgerUser {
	t.Helper()

	ledgerID, err := NewLedgerID()
	require.NoError(t, err)

	userID, err := entity.NewUserID()
	require.NoError(t, err)

	return NewLedgerUser(ledgerID, userID, RoleViewer)
}