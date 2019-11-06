package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestCreateRole(t *testing.T) {
	t.Parallel()

	t.Run("create_new_role_success", func(t *testing.T) {
		t.Parallel()

		roleName := testApp.mig.CreateUniqueString("role_name_")
		roleDesc := testApp.mig.CreateUniqueString("role_desc_")
		role, err := testApp.roleService.CreateRole(roleName, roleDesc)

		require.Nil(t, err)
		require.NotNil(t, role)
		require.Equal(t, role.RoleName, roleName)
		require.Equal(t, role.RoleDesc, roleDesc)
		require.True(t, role.Active)
	})

	t.Run("duplicated_role_name", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)

		role, err := testApp.roleService.CreateRole(role.RoleName, "vespa")

		require.Nil(t, role)
		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDuplicatedData))
	})
}

func TestGetRole(t *testing.T) {
	t.Parallel()

	t.Run("get_role_by_id_success", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)

		found, err := testApp.roleService.GetRole(role.ID)

		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.RoleName, role.RoleName)
		require.Equal(t, found.RoleDesc, role.RoleDesc)
		require.True(t, found.Active)
	})

	t.Run("get_role_by_id_not_found", func(t *testing.T) {
		t.Parallel()

		role, err := testApp.roleService.GetRole(int64(-1))

		require.Nil(t, err)
		require.Nil(t, role)
	})
}

func TestUpdateRole(t *testing.T) {
	t.Parallel()

	t.Run("update_role_by_success", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)

		updatedRole, err := testApp.roleService.UpdateRole(role.ID, "create_role_updated",
			"Create a role updated", nil)

		require.Nil(t, err)
		require.NotNil(t, updatedRole)
		require.Equal(t, updatedRole.RoleName, "create_role_updated")
		require.Equal(t, updatedRole.RoleDesc, "Create a role updated")
		require.True(t, updatedRole.Active)
	})

	t.Run("deactivate_role_uccess", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)
		activeStatus := false

		updatedRole, err := testApp.roleService.UpdateRole(role.ID, "", "", &activeStatus)

		require.Nil(t, err)
		require.NotNil(t, updatedRole)
		require.False(t, updatedRole.Active)
	})
}

func TestFetchRoles(t *testing.T) {
	t.Parallel()

	t.Run("fetch_roles_success", func(t *testing.T) {
		t.Parallel()

		testApp.mig.CreateSeedingRole(nil)
		testApp.mig.CreateSeedingRole(nil)

		roles, err := testApp.roleService.FetchRoles(2, "", nil, "+role_name")

		require.Nil(t, err)
		require.NotNil(t, roles)
		require.Len(t, roles, 2)
	})
}
