package mysqlrepo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryAction(t *testing.T) {
	t.Parallel()

	ids, err := testApp.loadActionFixtures("action_fix")
	if err != nil {
		log.Fatal(err)
		return
	}

	t.Run("query_action_by_name_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"action_name": "action_fix"}
		sort := map[string]comtype.SortDirection{"action_name": comtype.Ascending}
		actions, total, err := testApp.actionRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, actions)
		require.Len(t, actions, 10)
		require.Equal(t, int64(20), total)
	})

	t.Run("query_action_by_active_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"active": false}
		sort := map[string]comtype.SortDirection{}
		actions, total, err := testApp.actionRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, actions)
		require.Len(t, actions, 0)
		require.Zero(t, total)
	})

	t.Run("get_action_by_id_success", func(t *testing.T) {
		t.Parallel()

		action, err := testApp.actionRepo.GetByID(ids[0])

		require.Nil(t, err)
		require.NotNil(t, action)
		require.Equal(t, action.ID, ids[0])
	})

	t.Run("get_action_by_name_success", func(t *testing.T) {
		t.Parallel()

		name := testApp.generateUniqueString("test_action")
		id, err := testApp.createActionWithName(name)
		require.Nil(t, err)
		require.NotZero(t, id)

		action, err := testApp.actionRepo.GetByName(name)
		require.Nil(t, err)
		require.NotNil(t, action)
		require.Equal(t, id, action.ID)
		require.True(t, action.Active)
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

		actionName := testApp.generateUniqueString("created_action")
		id, err := testApp.actionRepo.Create(actionName, "created_action_desc")
		require.Nil(t, err)
		require.NotZero(t, id)
	})
}
