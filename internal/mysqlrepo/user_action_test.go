package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryUserAction(t *testing.T) {
	t.Parallel()

	t.Run("query_user_action_by_id_success", func(t *testing.T) {
		t.Parallel()

		ua, err := testApp.userActionRepo.GetByID(1)

		require.Nil(t, err)
		require.NotNil(t, ua)
		require.NotNil(t, ua.User)
		require.NotNil(t, ua.Action)
		require.Equal(t, int64(1), ua.ID)
		require.Equal(t, int64(2), ua.User.ID)
		require.Equal(t, int64(10), ua.Action.ID)
	})

	t.Run("query_user_action_by_id_not_found", func(t *testing.T) {
		t.Parallel()

		ua, err := testApp.userActionRepo.GetByID(-1)

		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDataNotFound))
		require.Nil(t, ua)
	})

	t.Run("query_user_action_by_user_id_success", func(t *testing.T) {
		t.Parallel()

		uas, err := testApp.userActionRepo.Query(12, map[string]interface{}{
			"user_id": 2,
		})

		require.Nil(t, err)
		require.NotNil(t, uas)
		require.Len(t, uas, 1)
		require.NotNil(t, uas[0].User)
		require.Equal(t, int64(2), uas[0].User.ID)
		require.Equal(t, int64(10), uas[0].Action.ID)
	})

	t.Run("query_user_action_by_action_id_fail", func(t *testing.T) {
		t.Parallel()

		uas, err := testApp.userActionRepo.Query(12, map[string]interface{}{
			"action_id": 999999999,
		})

		require.Nil(t, err)
		require.Len(t, uas, 0)
	})
}

func TestCreateUserAction(t *testing.T) {
	t.Parallel()

	t.Run("create_delete_user_action_success", func(t *testing.T) {
		t.Parallel()

		id, err := testApp.userActionRepo.Create(int64(1), int64(9))
		require.Nil(t, err)
		require.NotZero(t, id)

		err = testApp.userActionRepo.Delete(id)
		require.Nil(t, err)
	})
}
