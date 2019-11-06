package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestCreateRoleAction(t *testing.T) {
	t.Parallel()

	t.Run("create_role_action_success", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)
		action := testApp.mig.CreateSeedingAction(nil)
		roleAction, err := testApp.roleActionService.CreateRoleAction(role.ID, action.ID)

		require.Nil(t, err)
		require.NotNil(t, roleAction)
		require.Equal(t, roleAction.Role.ID, role.ID)
		require.Equal(t, roleAction.Action.ID, action.ID)
	})

	t.Run("create_role_action_fail", func(t *testing.T) {
		t.Parallel()

		roleAction, err := testApp.roleActionService.CreateRoleAction(int64(999999999), int64(999999999))

		require.NotNil(t, err)
		require.Nil(t, roleAction)
		require.True(t, err.Is(comtype.ErrDataNotFound))
	})
}

func TestDeleteRoleAction(t *testing.T) {
	t.Parallel()

	t.Run("delete_role_action_success", func(t *testing.T) {
		t.Parallel()

		roleAction := testApp.mig.CreateSeedingRoleAction(nil)
		err := testApp.roleActionService.DeleteRoleAction(roleAction.RoleID, roleAction.ActionID)

		require.Nil(t, err)
	})
}

func TestGetRoleAction(t *testing.T) {
	t.Parallel()

	t.Run("get_role_action_success", func(t *testing.T) {
		t.Parallel()

		roleAction := testApp.mig.CreateSeedingRoleAction(nil)
		found, err := testApp.roleActionService.GetRoleAction(roleAction.ID)

		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.ID, roleAction.ID)
		require.Equal(t, found.RoleID, roleAction.RoleID)
		require.Equal(t, found.ActionID, roleAction.ActionID)
	})
}

func TestQueryRoleAction(t *testing.T) {
	t.Parallel()

	t.Run("query_role_actions_role_id", func(t *testing.T) {
		t.Parallel()

		role := testApp.mig.CreateSeedingRole(nil)
		action1 := testApp.mig.CreateSeedingAction(nil)
		action2 := testApp.mig.CreateSeedingAction(nil)

		testApp.mig.CreateSeedingRoleAction(func(fields map[string]interface{}) {
			fields["role_id"] = role.ID
			fields["action_id"] = action1.ID
		})
		testApp.mig.CreateSeedingRoleAction(func(fields map[string]interface{}) {
			fields["role_id"] = role.ID
			fields["action_id"] = action2.ID
		})

		roleActions, err := testApp.roleActionService.FetchRoleActions(3, role.ID, 0)

		require.Nil(t, err)
		require.NotNil(t, roleActions)
		require.Len(t, roleActions, 2)
	})

	t.Run("query_role_actions_action_id", func(t *testing.T) {
		t.Parallel()

		action := testApp.mig.CreateSeedingAction(nil)
		role1 := testApp.mig.CreateSeedingRole(nil)
		role2 := testApp.mig.CreateSeedingRole(nil)

		testApp.mig.CreateSeedingRoleAction(func(fields map[string]interface{}) {
			fields["role_id"] = role1.ID
			fields["action_id"] = action.ID
		})
		testApp.mig.CreateSeedingRoleAction(func(fields map[string]interface{}) {
			fields["role_id"] = role2.ID
			fields["action_id"] = action.ID
		})

		roleActions, err := testApp.roleActionService.FetchRoleActions(3, 0, action.ID)

		require.Nil(t, err)
		require.NotNil(t, roleActions)
		require.Len(t, roleActions, 2)
	})
}
