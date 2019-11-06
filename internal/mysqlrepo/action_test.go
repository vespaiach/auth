package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryAction(t *testing.T) {
	t.Parallel()

	t.Run("query_action_by_name_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"action_name": "_action"}
		sort := map[string]comtype.SortDirection{"action_name": comtype.Ascending}
		actions, err := testApp.actionRepo.Query(10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, actions)
		require.Greater(t, len(actions), 0)
	})

	t.Run("query_action_by_active_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"active": false}
		sort := map[string]comtype.SortDirection{}
		actions, err := testApp.actionRepo.Query(10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, actions)
		require.Len(t, actions, 0)
	})

	t.Run("get_action_by_id_success", func(t *testing.T) {
		t.Parallel()

		action, err := testApp.actionRepo.GetByID(1)

		require.Nil(t, err)
		require.NotNil(t, action)
		require.Equal(t, action.ID, int64(1))
		require.Equal(t, action.ActionName, "create_action")
		require.Equal(t, action.ActionDesc, "Create a action")
	})

	t.Run("get_action_by_id_not_found", func(t *testing.T) {
		t.Parallel()

		action, err := testApp.actionRepo.GetByID(-1)

		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDataNotFound))
		require.Nil(t, action)
	})

	t.Run("get_action_by_name_success", func(t *testing.T) {
		t.Parallel()

		action, err := testApp.actionRepo.GetByName("list_user")
		require.Nil(t, err)
		require.NotNil(t, action)
		require.Equal(t, int64(10), action.ID)
		require.True(t, action.Active)
	})

	t.Run("get_action_by_user_id_success", func(t *testing.T) {
		t.Parallel()

		actions, err := testApp.actionRepo.GetByUserID(int64(1))
		require.Nil(t, err)
		require.NotNil(t, actions)
		require.Greater(t, len(actions), 0)
	})
}

func TestCreateAction(t *testing.T) {
	t.Parallel()

	t.Run("create_action_success", func(t *testing.T) {
		t.Parallel()

		actionName := testApp.generateUniqueString("created_action")

		id, err := testApp.actionRepo.Create(actionName, "created_action_desc")
		require.Nil(t, err)
		require.NotZero(t, id)

		found, err := testApp.actionRepo.GetByID(id)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, id, found.ID)
		require.True(t, found.Active)
		require.Equal(t, actionName, found.ActionName)
		require.Equal(t, "created_action_desc", found.ActionDesc)
	})

	t.Run("create_action_fail", func(t *testing.T) {
		t.Parallel()

		id, err := testApp.actionRepo.Create("create_action", "create_action")
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}

func TestUpdateAction(t *testing.T) {
	t.Parallel()

	t.Run("update_action_success", func(t *testing.T) {
		t.Parallel()

		action := testApp.mig.CreateSeedingAction(nil)

		fields := map[string]interface{}{
			"action_name": "updated",
			"action_desc": "updated",
			"active":      comtype.Active,
		}
		err := testApp.actionRepo.Update(action.ID, fields)
		require.Nil(t, err)

		found, err := testApp.actionRepo.GetByID(action.ID)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.ActionName, "updated")
		require.Equal(t, found.ActionDesc, "updated")
		require.True(t, found.Active)
	})

	t.Run("update_duplicated_action_name", func(t *testing.T) {
		t.Parallel()

		action1 := testApp.mig.CreateSeedingAction(nil)
		action2 := testApp.mig.CreateSeedingAction(nil)

		fields := map[string]interface{}{
			"action_name": action1.ActionName,
		}
		err := testApp.actionRepo.Update(action2.ID, fields)
		require.NotNil(t, err)
	})
}
