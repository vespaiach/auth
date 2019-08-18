package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryUser(t *testing.T) {
	t.Parallel()

	t.Run("query_user_by_name_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"username": "admin"}
		sort := map[string]comtype.SortDirection{"full_name": comtype.Ascending}
		users, total, err := testApp.userRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, users)
		require.Len(t, users, 1)
		require.Equal(t, int64(1), total)
	})

	t.Run("query_user_by_active_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"active": false}
		sort := map[string]comtype.SortDirection{}
		users, total, err := testApp.userRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, users)
		require.Len(t, users, 0)
		require.Zero(t, total)
	})

	t.Run("get_user_by_id_success", func(t *testing.T) {
		t.Parallel()

		user, err := testApp.userRepo.GetByID(1)

		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, int64(1), user.ID)
		require.Equal(t, "admin", user.Username)
		require.Equal(t, "admin@test.com", user.Email)
	})

	t.Run("get_user_by_username_success", func(t *testing.T) {
		t.Parallel()

		user, err := testApp.userRepo.GetByUsername("staff")
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, int64(2), user.ID)
		require.Equal(t, "staff", user.Username)
		require.True(t, user.Active)
	})
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("create_user_success", func(t *testing.T) {
		t.Parallel()

		userName := testApp.generateUniqueString("created_user")

		id, err := testApp.userRepo.Create("full_name", userName, "password", "email")
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("create_user_fail", func(t *testing.T) {
		t.Parallel()

		id, err := testApp.userRepo.Create("full_name_admin", "admin", "password", "email")
		require.NotNil(t, err)
		require.Zero(t, id)

		id, err = testApp.userRepo.Create("full_name_admin", "admin_test", "password", "admin@test.com")
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}
