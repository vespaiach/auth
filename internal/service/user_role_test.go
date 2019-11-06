package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestCreateUserRole(t *testing.T) {
	t.Parallel()

	t.Run("create_user_role_success", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)
		user := testApp.mig.CreateSeedingUser(nil)
		userRole, err := testApp.userRoleService.CreateUserRole(user.ID, role.ID)

		require.Nil(t, err)
		require.NotNil(t, userRole)
		require.Equal(t, userRole.Role.ID, role.ID)
		require.Equal(t, userRole.User.ID, user.ID)
	})

	t.Run("create_user_role_fail", func(t *testing.T) {
		t.Parallel()

		roleAction, err := testApp.userRoleService.CreateUserRole(int64(999999999), int64(999999999))

		require.NotNil(t, err)
		require.Nil(t, roleAction)
		require.True(t, err.Is(comtype.ErrDataNotFound))
	})
}

func TestDeleteUserRole(t *testing.T) {
	t.Parallel()

	t.Run("delete_user_role_success", func(t *testing.T) {
		t.Parallel()

		userRole := testApp.mig.CreateSeedingUserRole(nil)
		err := testApp.userRoleService.DeleteUserRole(userRole.UserID, userRole.RoleID)

		require.Nil(t, err)
	})
}

func TestGetUserRole(t *testing.T) {
	t.Parallel()

	t.Run("get_user_role_success", func(t *testing.T) {
		t.Parallel()

		userRole := testApp.mig.CreateSeedingUserRole(nil)
		found, err := testApp.userRoleService.GetUserRole(userRole.ID)

		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.ID, userRole.ID)
		require.Equal(t, found.User.ID, userRole.UserID)
		require.Equal(t, found.Role.ID, userRole.RoleID)
	})
}

func TestQueryUserRole(t *testing.T) {
	t.Parallel()

	t.Run("query_user_roles_role_id", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)
		user1 := testApp.mig.CreateSeedingUser(nil)
		user2 := testApp.mig.CreateSeedingUser(nil)

		testApp.mig.CreateSeedingUserRole(func(fields map[string]interface{}) {
			fields["role_id"] = role.ID
			fields["user_id"] = user1.ID
		})
		testApp.mig.CreateSeedingUserRole(func(fields map[string]interface{}) {
			fields["role_id"] = role.ID
			fields["user_id"] = user2.ID
		})

		userRoles, err := testApp.userRoleService.FetchUserRoles(3, 0, role.ID)

		require.Nil(t, err)
		require.NotNil(t, userRoles)
		require.Len(t, userRoles, 2)
	})

	t.Run("query_user_roles_user_id", func(t *testing.T) {
		t.Parallel()

		user := testApp.mig.CreateSeedingUser(nil)
		role1 := testApp.mig.CreateSeedingRole(nil)
		role2 := testApp.mig.CreateSeedingRole(nil)

		testApp.mig.CreateSeedingUserRole(func(fields map[string]interface{}) {
			fields["role_id"] = role1.ID
			fields["user_id"] = user.ID
		})
		testApp.mig.CreateSeedingUserRole(func(fields map[string]interface{}) {
			fields["role_id"] = role2.ID
			fields["user_id"] = user.ID
		})

		userRoles, err := testApp.userRoleService.FetchUserRoles(3, user.ID, 0)

		require.Nil(t, err)
		require.NotNil(t, userRoles)
		require.Len(t, userRoles, 2)
	})
}
