package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryRoleAction(t *testing.T) {
	t.Parallel()

	t.Run("query_role_action_by_id_success", func(t *testing.T) {
		t.Parallel()

		ra, err := testApp.roleActionRepo.GetByID(1)

		require.Nil(t, err)
		require.NotNil(t, ra)
		require.NotNil(t, ra.Role)
		require.NotNil(t, ra.Action)
		require.Equal(t, int64(1), ra.ID)
		require.Equal(t, int64(1), ra.Role.ID)
		require.Equal(t, int64(1), ra.Action.ID)
	})

	t.Run("query_role_action_by_id_not_found", func(t *testing.T) {
		t.Parallel()

		ra, err := testApp.roleActionRepo.GetByID(-1)

		require.NotNil(t, err)
		require.NotNil(t, err.Is(comtype.ErrDataNotFound))
		require.Nil(t, ra)
	})

	t.Run("query_role_action_by_role_id_success", func(t *testing.T) {
		t.Parallel()

		ras, err := testApp.roleActionRepo.Query(12, map[string]interface{}{
			"role_id": 1,
		})

		require.Nil(t, err)
		require.NotNil(t, ras)
		require.Greater(t, len(ras), 0)
		require.NotNil(t, ras[0].Role)
		require.NotNil(t, ras[0].Action)
	})

	t.Run("query_role_action_by_action_id_empty", func(t *testing.T) {
		t.Parallel()

		ras, err := testApp.roleActionRepo.Query(12, map[string]interface{}{
			"action_id": 999999999,
		})

		require.Nil(t, err)
		require.NotNil(t, ras)
		require.Len(t, ras, 0)
	})
}

func TestCreateRoleAction(t *testing.T) {
	t.Parallel()

	t.Run("create_delete_role_action_success", func(t *testing.T) {
		t.Parallel()

		id, err := testApp.roleActionRepo.Create(int64(2), int64(6))
		require.Nil(t, err)
		require.NotZero(t, id)

		err = testApp.roleActionRepo.Delete(id)
		require.Nil(t, err)
	})
}
