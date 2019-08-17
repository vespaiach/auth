package mysqlrepo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryRole(t *testing.T) {
	t.Parallel()

	ids, err := testApp.loadRoleFixtures("role_fixs")
	if err != nil {
		log.Fatal(err)
		return
	}

	t.Run("query_role_by_name_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"role_name": "role_fixs"}
		sort := map[string]comtype.SortDirection{"role_name": comtype.Ascending}
		roles, total, err := testApp.roleRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, roles)
		require.Len(t, roles, 10)
		require.Equal(t, int64(20), total)
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

		role, err := testApp.roleRepo.GetByID(ids[0])

		require.Nil(t, err)
		require.NotNil(t, role)
		require.Equal(t, role.ID, ids[0])
	})

	t.Run("get_role_by_name_success", func(t *testing.T) {
		t.Parallel()

		name := testApp.generateUniqueString("test_role")
		id, err := testApp.createRoleWithName(name)
		require.Nil(t, err)
		require.NotZero(t, id)

		role, err := testApp.roleRepo.GetByName(name)
		require.Nil(t, err)
		require.NotNil(t, role)
		require.Equal(t, id, role.ID)
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

		roleName := testApp.generateUniqueString("created_role")
		id, err := testApp.roleRepo.Create(roleName, "created_role_desc")
		require.Nil(t, err)
		require.NotZero(t, id)
	})
}
