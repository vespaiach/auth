package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestCreateAction(t *testing.T) {
	t.Parallel()

	t.Run("create_new_action_success", func(t *testing.T) {
		t.Parallel()

		actionName := testApp.mig.CreateUniqueString("action_name_")
		actionDesc := testApp.mig.CreateUniqueString("action_desc_")
		action, err := testApp.actionService.CreateAction(actionName, actionDesc)

		require.Nil(t, err)
		require.NotNil(t, action)
		require.Equal(t, action.ActionName, actionName)
		require.Equal(t, action.ActionDesc, actionDesc)
		require.True(t, action.Active)
	})

	t.Run("duplicated_action_name", func(t *testing.T) {
		t.Parallel()

		action := testApp.mig.CreateSeedingAction(nil)

		action, err := testApp.actionService.CreateAction(action.ActionName, "vespa")

		require.Nil(t, action)
		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDuplicatedData))
	})
}

func TestGetAction(t *testing.T) {
	t.Parallel()

	t.Run("get_action_by_id_success", func(t *testing.T) {
		t.Parallel()

		action, err := testApp.actionService.GetAction(int64(1))

		require.Nil(t, err)
		require.NotNil(t, action)
		require.Equal(t, action.ActionName, "create_action")
		require.Equal(t, action.ActionDesc, "Create a action")
		require.True(t, action.Active)
	})

	t.Run("get_action_by_id_not_found", func(t *testing.T) {
		t.Parallel()

		action, err := testApp.actionService.GetAction(int64(-1))

		require.Nil(t, err)
		require.Nil(t, action)
	})
}

func TestUpdateAction(t *testing.T) {
	t.Parallel()

	t.Run("update_action_by_success", func(t *testing.T) {
		t.Parallel()

		action := testApp.mig.CreateSeedingAction(nil)

		updatedAction, err := testApp.actionService.UpdateAction(action.ID, "create_action_updated",
			"Create a action updated", nil)

		require.Nil(t, err)
		require.NotNil(t, updatedAction)
		require.Equal(t, updatedAction.ActionName, "create_action_updated")
		require.Equal(t, updatedAction.ActionDesc, "Create a action updated")
		require.True(t, updatedAction.Active)
	})
}

func TestFetchActions(t *testing.T) {
	t.Parallel()

	t.Run("fetch_actions_success", func(t *testing.T) {
		t.Parallel()

		actions, err := testApp.actionService.FetchActions(5, "", nil, "+action_name")

		require.Nil(t, err)
		require.NotNil(t, actions)
		require.Len(t, actions, 5)
	})
}
