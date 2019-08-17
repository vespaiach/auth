package mysqlrepo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestQueryUser(t *testing.T) {
	t.Parallel()

	ids, err := testApp.loadUserFixtures("user_fixs")
	if err != nil {
		log.Fatal(err)
		return
	}

	t.Run("query_user_by_name_success", func(t *testing.T) {
		t.Parallel()

		filter := map[string]interface{}{"username": "user_fixs"}
		sort := map[string]comtype.SortDirection{"full_name": comtype.Ascending}
		users, total, err := testApp.userRepo.Query(1, 10, filter, sort)

		require.Nil(t, err)
		require.NotNil(t, users)
		require.Len(t, users, 10)
		require.Equal(t, int64(20), total)
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

		user, err := testApp.userRepo.GetByID(ids[0])

		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.ID, ids[0])
	})

	t.Run("get_user_by_username_success", func(t *testing.T) {
		t.Parallel()

		name := testApp.generateUniqueString("test_user")
		id, err := testApp.createUserWithName(name)
		require.Nil(t, err)
		require.NotZero(t, id)

		user, err := testApp.userRepo.GetByUsername(name)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, id, user.ID)
		require.True(t, user.Active)
	})
}

// func TestCreateUser(t *testing.T) {
// 	t.Parallel()

// 	t.Run("create_user_success", func(t *testing.T) {
// 		t.Parallel()

// 		userName := testApp.generateUniqueString("created_user")

// 		id, err := testApp.userRepo.Create(userName, "created_user_desc")
// 		require.Nil(t, err)
// 		require.NotZero(t, id)

// 		found, err := testApp.userRepo.GetByID(id)
// 		require.Nil(t, err)
// 		require.NotNil(t, found)
// 		require.Equal(t, id, found.ID)
// 		require.True(t, found.Active)
// 		require.Equal(t, userName, found.UserName)
// 		require.Equal(t, "created_user_desc", found.UserDesc)
// 	})

// 	t.Run("create_user_fail", func(t *testing.T) {
// 		t.Parallel()

// 		userName := testApp.generateUniqueString("created_user")
// 		id, err := testApp.userRepo.Create(userName, "created_user_desc")
// 		require.Nil(t, err)
// 		require.NotZero(t, id)
// 	})
// }
