package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryRole(t *testing.T) {
	t.Parallel()

	t.Run("query_role_by_name_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"role_name": "_role"}
		sort := map[string]comtype.SortDirection{"role_name": comtype.Ascending}
		roles, total, err := testApp.roleRepo.Query(1, 2, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, roles)
		require.Greater(t, len(roles), 0)
		require.NotZero(t, total)
	})

	t.Run("query_role_by_active_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"active": false}
		sort := map[string]comtype.SortDirection{}
		roles, total, err := testApp.roleRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, roles)
		require.Len(t, roles, 0)
		require.Zero(t, total)
	})

	t.Run("get_role_by_id_success", func(t *testing.T) {
		t.Parallel()

		role, err := testApp.roleRepo.GetByID(1)

		require.Nil(t, err)
		require.NotNil(t, role)
		require.Equal(t, role.ID, int64(1))
		require.Equal(t, role.RoleName, "admin_role")
		require.Equal(t, role.RoleDesc, "Admin role")
	})

	t.Run("get_role_by_name_success", func(t *testing.T) {
		t.Parallel()

		role, err := testApp.roleRepo.GetByName("staff_role")
		require.Nil(t, err)
		require.NotNil(t, role)
		require.Equal(t, int64(2), role.ID)
		require.Equal(t, "staff_role", role.RoleName)
		require.True(t, role.Active)
	})
}

func TestCreateRole(t *testing.T) {
	t.Parallel()

	t.Run("create_role_success", func(t *testing.T) {
		t.Parallel()

		roleName := testApp.generateUniqueString("created_role")

		id, err := testApp.roleRepo.Create(roleName, "created_role_desc")
		require.Nil(t, err)
		require.NotZero(t, id)

		found, err := testApp.roleRepo.GetByID(id)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, id, found.ID)
		require.True(t, found.Active)
		require.Equal(t, roleName, found.RoleName)
		require.Equal(t, "created_role_desc", found.RoleDesc)
	})

	t.Run("create_role_fail", func(t *testing.T) {
		t.Parallel()

		id, err := testApp.roleRepo.Create("admin_role", "created_role_desc")
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}
