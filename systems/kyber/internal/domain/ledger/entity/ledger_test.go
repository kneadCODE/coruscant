package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kneadCODE/coruscant/systems/kyber/internal/domain/user/entity"
	"github.com/kneadCODE/coruscant/systems/kyber/internal/pkg/money"
)

func TestNewLedger(t *testing.T) {
	adminUserID, err := entity.NewUserID()
	require.NoError(t, err)

	tests := []struct {
		name         string
		ledgerName   string
		description  string
		baseCurrency money.Currency
		adminUserID  entity.UserID
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid ledger creation",
			ledgerName:   "Personal Finances",
			description:  "My personal finance tracking",
			baseCurrency: "USD",
			adminUserID:  adminUserID,
			wantErr:      false,
		},
		{
			name:         "empty ledger name",
			ledgerName:   "",
			description:  "Test description",
			baseCurrency: "USD",
			adminUserID:  adminUserID,
			wantErr:      true,
			errContains:  "ledger name cannot be empty",
		},
		{
			name:         "empty base currency",
			ledgerName:   "Test Ledger",
			description:  "Test description",
			baseCurrency: "",
			adminUserID:  adminUserID,
			wantErr:      true,
			errContains:  "base currency cannot be empty",
		},
		{
			name:         "invalid admin user ID",
			ledgerName:   "Test Ledger",
			description:  "Test description",
			baseCurrency: "USD",
			adminUserID:  entity.UserID{},
			wantErr:      true,
			errContains:  "admin user ID is invalid",
		},
		{
			name:         "empty description is allowed",
			ledgerName:   "Test Ledger",
			description:  "",
			baseCurrency: "EUR",
			adminUserID:  adminUserID,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ledger, err := NewLedger(tt.ledgerName, tt.description, tt.baseCurrency, tt.adminUserID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, ledger)
			} else {
				require.NoError(t, err)
				require.NotNil(t, ledger)

				assert.True(t, ledger.ID.IsValid())
				assert.Equal(t, tt.ledgerName, ledger.Name)
				assert.Equal(t, tt.description, ledger.Description)
				assert.Equal(t, tt.baseCurrency, ledger.BaseCurrency)
				assert.Equal(t, LedgerStatusActive, ledger.Status)
				assert.Len(t, ledger.Users, 1)
				assert.Equal(t, tt.adminUserID, ledger.Users[0].UserID)
				assert.Equal(t, RoleAdmin, ledger.Users[0].Role)
				assert.False(t, ledger.CreatedAt.IsZero())
				assert.False(t, ledger.UpdatedAt.IsZero())
			}
		})
	}
}

func TestReconstructLedger(t *testing.T) {
	ledgerID, err := NewLedgerID()
	require.NoError(t, err)

	userID, err := entity.NewUserID()
	require.NoError(t, err)

	users := []LedgerUser{
		*NewLedgerUser(ledgerID, userID, RoleAdmin),
	}

	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	ledger := ReconstructLedger(
		ledgerID,
		"Reconstructed Ledger",
		"Test description",
		"USD",
		LedgerStatusArchived,
		users,
		createdAt,
		updatedAt,
	)

	assert.Equal(t, ledgerID, ledger.ID)
	assert.Equal(t, "Reconstructed Ledger", ledger.Name)
	assert.Equal(t, "Test description", ledger.Description)
	assert.Equal(t, money.Currency("USD"), ledger.BaseCurrency)
	assert.Equal(t, LedgerStatusArchived, ledger.Status)
	assert.Equal(t, users, ledger.Users)
	assert.Equal(t, createdAt, ledger.CreatedAt)
	assert.Equal(t, updatedAt, ledger.UpdatedAt)
}

func TestLedger_UpdateInfo(t *testing.T) {
	ledger := createTestLedger(t)
	originalUpdatedAt := ledger.UpdatedAt

	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		newName     string
		description string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid update",
			newName:     "Updated Ledger Name",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "empty name",
			newName:     "",
			description: "Some description",
			wantErr:     true,
			errContains: "ledger name cannot be empty",
		},
		{
			name:        "empty description is allowed",
			newName:     "Ledger Name",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ledger.UpdateInfo(tt.newName, tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newName, ledger.Name)
				assert.Equal(t, tt.description, ledger.Description)
				assert.True(t, ledger.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}

func TestLedger_UpdateUserRole(t *testing.T) {
	ledger := createTestLedger(t)

	// Add a second user
	viewerUserID, err := entity.NewUserID()
	require.NoError(t, err)

	ledgerUser := NewLedgerUser(ledger.ID, viewerUserID, RoleViewer)
	ledger.Users = append(ledger.Users, *ledgerUser)

	originalUpdatedAt := ledger.UpdatedAt
	time.Sleep(time.Millisecond)

	tests := []struct {
		name        string
		userID      entity.UserID
		newRole     Role
		wantErr     bool
		errContains string
	}{
		{
			name:    "update existing user role",
			userID:  viewerUserID,
			newRole: RoleEditor,
			wantErr: false,
		},
		{
			name:        "user not found",
			userID:      entity.UserID{}, // Invalid user ID
			newRole:     RoleViewer,
			wantErr:     true,
			errContains: "user not found in ledger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ledger.UpdateUserRole(tt.userID, tt.newRole)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.True(t, ledger.UpdatedAt.After(originalUpdatedAt))

				// Verify the role was updated
				userAccess, err := ledger.GetUserAccess(tt.userID)
				require.NoError(t, err)
				assert.Equal(t, tt.newRole, userAccess.Role)
			}
		})
	}
}

func TestLedger_GetUserAccess(t *testing.T) {
	ledger := createTestLedger(t)

	adminUserID := ledger.Users[0].UserID

	// Add another user
	editorUserID, err := entity.NewUserID()
	require.NoError(t, err)

	editorUser := NewLedgerUser(ledger.ID, editorUserID, RoleEditor)
	ledger.Users = append(ledger.Users, *editorUser)

	tests := []struct {
		name        string
		userID      entity.UserID
		expectFound bool
		expectedRole Role
	}{
		{
			name:         "get admin user access",
			userID:       adminUserID,
			expectFound:  true,
			expectedRole: RoleAdmin,
		},
		{
			name:         "get editor user access",
			userID:       editorUserID,
			expectFound:  true,
			expectedRole: RoleEditor,
		},
		{
			name:        "user not found",
			userID:      entity.UserID{},
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userAccess, err := ledger.GetUserAccess(tt.userID)

			if tt.expectFound {
				require.NoError(t, err)
				require.NotNil(t, userAccess)
				assert.Equal(t, tt.userID, userAccess.UserID)
				assert.Equal(t, tt.expectedRole, userAccess.Role)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user does not have access")
				assert.Nil(t, userAccess)
			}
		})
	}
}

func TestLedger_HasUserAccess(t *testing.T) {
	ledger := createTestLedger(t)
	adminUserID := ledger.Users[0].UserID

	// Add another user
	viewerUserID, err := entity.NewUserID()
	require.NoError(t, err)

	viewerUser := NewLedgerUser(ledger.ID, viewerUserID, RoleViewer)
	ledger.Users = append(ledger.Users, *viewerUser)

	tests := []struct {
		name     string
		userID   entity.UserID
		expected bool
	}{
		{
			name:     "admin has access",
			userID:   adminUserID,
			expected: true,
		},
		{
			name:     "viewer has access",
			userID:   viewerUserID,
			expected: true,
		},
		{
			name:     "unknown user has no access",
			userID:   entity.UserID{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ledger.HasUserAccess(tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLedger_UserHasPermission(t *testing.T) {
	ledger := createTestLedger(t)
	adminUserID := ledger.Users[0].UserID

	// Add users with different roles
	editorUserID, err := entity.NewUserID()
	require.NoError(t, err)
	ledger.Users = append(ledger.Users, *NewLedgerUser(ledger.ID, editorUserID, RoleEditor))

	viewerUserID, err := entity.NewUserID()
	require.NoError(t, err)
	ledger.Users = append(ledger.Users, *NewLedgerUser(ledger.ID, viewerUserID, RoleViewer))

	tests := []struct {
		name       string
		userID     entity.UserID
		permission Permission
		expected   bool
	}{
		{
			name:       "admin has admin permission",
			userID:     adminUserID,
			permission: PermissionAdmin,
			expected:   true,
		},
		{
			name:       "admin has edit permission",
			userID:     adminUserID,
			permission: PermissionEdit,
			expected:   true,
		},
		{
			name:       "admin has read permission",
			userID:     adminUserID,
			permission: PermissionReadOnly,
			expected:   true,
		},
		{
			name:       "editor has edit permission",
			userID:     editorUserID,
			permission: PermissionEdit,
			expected:   true,
		},
		{
			name:       "editor does not have admin permission",
			userID:     editorUserID,
			permission: PermissionAdmin,
			expected:   false,
		},
		{
			name:       "viewer has read permission",
			userID:     viewerUserID,
			permission: PermissionReadOnly,
			expected:   true,
		},
		{
			name:       "viewer does not have edit permission",
			userID:     viewerUserID,
			permission: PermissionEdit,
			expected:   false,
		},
		{
			name:       "unknown user has no permissions",
			userID:     entity.UserID{},
			permission: PermissionReadOnly,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ledger.UserHasPermission(tt.userID, tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLedger_GetAdmin(t *testing.T) {
	ledger := createTestLedger(t)
	originalAdminUserID := ledger.Users[0].UserID

	// Add a non-admin user
	viewerUserID, err := entity.NewUserID()
	require.NoError(t, err)
	ledger.Users = append(ledger.Users, *NewLedgerUser(ledger.ID, viewerUserID, RoleViewer))

	admin := ledger.GetAdmin()
	require.NotNil(t, admin)
	assert.Equal(t, originalAdminUserID, admin.UserID)
	assert.Equal(t, RoleAdmin, admin.Role)

	// Test ledger with no admin (edge case)
	ledger.Users = []LedgerUser{
		*NewLedgerUser(ledger.ID, viewerUserID, RoleViewer),
	}

	admin = ledger.GetAdmin()
	assert.Nil(t, admin)
}

func TestLedger_Archive(t *testing.T) {
	ledger := createTestLedger(t)
	assert.Equal(t, LedgerStatusActive, ledger.Status)

	originalUpdatedAt := ledger.UpdatedAt
	time.Sleep(time.Millisecond)

	err := ledger.Archive()
	require.NoError(t, err)
	assert.Equal(t, LedgerStatusArchived, ledger.Status)
	assert.True(t, ledger.UpdatedAt.After(originalUpdatedAt))

	// Test archiving an already archived ledger
	err = ledger.Archive()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ledger is already archived")
}

func TestLedger_Activate(t *testing.T) {
	ledger := createTestLedger(t)

	// First archive it
	err := ledger.Archive()
	require.NoError(t, err)
	assert.Equal(t, LedgerStatusArchived, ledger.Status)

	originalUpdatedAt := ledger.UpdatedAt
	time.Sleep(time.Millisecond)

	// Now activate it
	err = ledger.Activate()
	require.NoError(t, err)
	assert.Equal(t, LedgerStatusActive, ledger.Status)
	assert.True(t, ledger.UpdatedAt.After(originalUpdatedAt))

	// Test activating an already active ledger
	err = ledger.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ledger is already active")
}

func TestLedger_CanWrite(t *testing.T) {
	ledger := createTestLedger(t)

	// Active ledger can write
	assert.True(t, ledger.CanWrite())

	// Archived ledger cannot write
	err := ledger.Archive()
	require.NoError(t, err)
	assert.False(t, ledger.CanWrite())
}

func TestLedger_CanRead(t *testing.T) {
	ledger := createTestLedger(t)

	// Active ledger can read
	assert.True(t, ledger.CanRead())

	// Archived ledger can still read
	err := ledger.Archive()
	require.NoError(t, err)
	assert.True(t, ledger.CanRead())
}

func TestLedger_UserCount(t *testing.T) {
	ledger := createTestLedger(t)
	assert.Equal(t, 1, ledger.UserCount())

	// Add more users
	userID1, err := entity.NewUserID()
	require.NoError(t, err)
	ledger.Users = append(ledger.Users, *NewLedgerUser(ledger.ID, userID1, RoleViewer))

	userID2, err := entity.NewUserID()
	require.NoError(t, err)
	ledger.Users = append(ledger.Users, *NewLedgerUser(ledger.ID, userID2, RoleEditor))

	assert.Equal(t, 3, ledger.UserCount())
}

// Helper functions

func createTestLedger(t *testing.T) *Ledger {
	t.Helper()

	adminUserID, err := entity.NewUserID()
	require.NoError(t, err)

	ledger, err := NewLedger("Test Ledger", "Test description", "USD", adminUserID)
	require.NoError(t, err)

	return ledger
}