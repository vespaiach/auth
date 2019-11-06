package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryUserRole(t *testing.T) {
	t.Parallel()

	t.Run("query_user_role_by_id_success", func(t *testing.T) {
		t.Parallel()

		ur, err := testApp.userRoleRepo.GetByID(2)

		require.Nil(t, err)
		require.NotNil(t, ur)
		require.NotNil(t, ur.User)
		require.NotNil(t, ur.Role)
		require.Equal(t, int64(2), ur.ID)
		require.Equal(t, int64(1), ur.User.ID)
		require.Equal(t, int64(2), ur.Role.ID)
	})

	t.Run("query_user_role_by_id_not_found", func(t *testing.T) {
		t.Parallel()

		ur, err := testApp.userRoleRepo.GetByID(-2)

		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDataNotFound))
		require.Nil(t, ur)
	})

	t.Run("query_user_role_by_user_id_success", func(t *testing.T) {
		t.Parallel()

		urs, err := testApp.userRoleRepo.Query(12, map[string]interface{}{
			"user_id": 1,
		})

		require.Nil(t, err)
		require.NotNil(t, urs)
		require.Len(t, urs, 2)
		require.NotNil(t, urs[0].User)
		require.NotNil(t, urs[0].Role)
		require.NotNil(t, urs[1].User)
		require.NotNil(t, urs[1].Role)
	})

	t.Run("query_user_role_by_role_id_fail", func(t *testing.T) {
		t.Parallel()

		urs, err := testApp.userRoleRepo.Query(12, map[string]interface{}{
			"role_id": 999999999,
		})

		require.Nil(t, err)
		require.Len(t, urs, 0)
	})
}

func TestCreateUserRole(t *testing.T) {
	t.Parallel()

	t.Run("create_and_delete_user_role_success", func(t *testing.T) {
		t.Parallel()

		id, err := testApp.userRoleRepo.Create(int64(2), int64(1))
		require.Nil(t, err)
		require.NotZero(t, id)

		err = testApp.userRoleRepo.Delete(id)
		require.Nil(t, err)
	})
}
