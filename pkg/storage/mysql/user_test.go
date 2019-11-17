package mysql

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/common"
	"testing"
)

func TestUserStorage_AddUser(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_user", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")

		id, err := test.ust.AddUser(name, email, "hash")
		require.Nil(t, err)
		require.NotZero(t, id)
	})
}

func TestUserStorage_ModifyUser(t *testing.T) {
	t.Parallel()

	t.Run("success_modify_a_user", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingUser(nil)
		newname := test.mig.createUniqueString("username")
		newemail := test.mig.createUniqueString("email")

		err := test.ust.ModifyUser(id, newname, newemail, "hash_updated", sql.NullBool{Valid: true})
		require.Nil(t, err)

		name, email, hash, active := test.mig.getUserByID(id)
		require.Equal(t, name, newname)
		require.Equal(t, email, newemail)
		require.Equal(t, hash, "hash_updated")
		require.False(t, active)
	})
}

func TestUserStorage_GetUser(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_user_by_id", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingUser(nil)

		user, err := test.ust.GetUser(id)
		require.Nil(t, err)
		require.NotNil(t, user)
	})
}

func TestUserStorage_GetUserByEmail(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_user_by_email", func(t *testing.T) {
		t.Parallel()

		email := test.mig.createUniqueString("email")
		id := test.mig.createSeedingUser(func(field map[string]interface{}) {
			field["email"] = email
		})

		user, err := test.ust.GetUserByEmail(email)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.ID, id)
	})
}

func TestUserStorage_GetUserByUsername(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_user_by_username", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("username")
		id := test.mig.createSeedingUser(func(field map[string]interface{}) {
			field["username"] = name
		})

		user, err := test.ust.GetUserByUsername(name)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.ID, id)
	})
}

func TestUserStorage_QueryUsers(t *testing.T) {
	t.Parallel()

	t.Run("success_query_users", func(t *testing.T) {
		t.Parallel()

		name1 := test.mig.createUniqueString("user1ame")
		name2 := test.mig.createUniqueString("user1ame")
		name3 := test.mig.createUniqueString("user1ame")
		name4 := test.mig.createUniqueString("user1ame")
		name5 := test.mig.createUniqueString("user1ame")

		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name1 })
		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name2 })
		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name3 })
		test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = name4 })
		test.mig.createSeedingUser(func(field map[string]interface{}) {
			field["username"] = name5
			field["active"] = false
		})

		users, total, err := test.ust.QueryUsers(2, 2, "user1ame", "",
			sql.NullBool{Valid: true, Bool: true}, "username", common.Ascending)
		require.Nil(t, err)
		require.NotNil(t, users)
		require.Equal(t, int64(4), total)
		require.Len(t, users, 2)
	})
}

func TestUserStorage_GetBunches(t *testing.T) {
	t.Parallel()

	t.Run("success_get_bunches_of_user", func(t *testing.T) {
		t.Parallel()

		username := test.mig.createUniqueString("username")

		id1 := test.mig.createSeedingBunch(nil)
		id2 := test.mig.createSeedingBunch(nil)
		id3 := test.mig.createSeedingBunch(nil)
		id4 := test.mig.createSeedingBunch(nil)

		userID := test.mig.createSeedingUser(func(field map[string]interface{}) { field["username"] = username })

		err := test.ust.AddBunchesToUser(userID, []int64{id1, id2, id3, id4})
		require.Nil(t, err)

		bunches, err := test.ust.GetBunches(username)
		require.Nil(t, err)
		require.NotNil(t, bunches)
		require.Len(t, bunches, 4)
	})
}

func TestUserStorage_AddBunchesToUser(t *testing.T) {
	t.Parallel()

	t.Run("success_add_bunches_to_user", func(t *testing.T) {
		t.Parallel()

		id1 := test.mig.createSeedingBunch(nil)
		id2 := test.mig.createSeedingBunch(nil)
		id3 := test.mig.createSeedingBunch(nil)
		id4 := test.mig.createSeedingBunch(nil)

		userID := test.mig.createSeedingUser(nil)

		err := test.ust.AddBunchesToUser(userID, []int64{id1, id2, id3, id4})
		require.Nil(t, err)
	})
}

func TestUserStorage_GetBunchIDs(t *testing.T) {
	t.Parallel()

	t.Run("success_get_bunch_ids_by_name", func(t *testing.T) {
		t.Parallel()

		name1 := test.mig.createUniqueString("bunch")
		name2 := test.mig.createUniqueString("bunch")
		name3 := test.mig.createUniqueString("bunch")
		name4 := test.mig.createUniqueString("bunch")

		id1 := test.mig.createSeedingBunch(func(field map[string]interface{}) { field["name"] = name1 })
		id2 := test.mig.createSeedingBunch(func(field map[string]interface{}) { field["name"] = name2 })
		id3 := test.mig.createSeedingBunch(func(field map[string]interface{}) { field["name"] = name3 })
		id4 := test.mig.createSeedingBunch(func(field map[string]interface{}) { field["name"] = name4 })

		bIDs, err := test.ust.GetBunchIDs([]string{name1, name2, name3, name4})
		require.Nil(t, err)
		require.NotNil(t, bIDs)
		require.Len(t, bIDs, 4)
		require.True(t, contains(bIDs, id1))
		require.True(t, contains(bIDs, id2))
		require.True(t, contains(bIDs, id3))
		require.True(t, contains(bIDs, id4))
	})
}

func TestUserStorage_GetKeys(t *testing.T) {
	t.Parallel()

	t.Run("success_get_keys_ids_by_username", func(t *testing.T) {
		t.Parallel()

		username := test.mig.createUniqueString("username")
		kID1 := test.mig.createSeedingServiceKey(nil)
		kID2 := test.mig.createSeedingServiceKey(nil)
		kID3 := test.mig.createSeedingServiceKey(nil)
		kID4 := test.mig.createSeedingServiceKey(nil)
		bID1 := test.mig.createSeedingBunch(nil)
		bID2 := test.mig.createSeedingBunch(nil)
		uID := test.mig.createSeedingUser(func(field map[string]interface{}) {
			field["username"] = username
		})

		err := test.bst.AddKeysToBunch(bID1, []int64{kID1, kID2, kID3})
		require.Nil(t, err)

		err = test.bst.AddKeysToBunch(bID2, []int64{kID3, kID4})
		require.Nil(t, err)

		err = test.ust.AddBunchesToUser(uID, []int64{bID1, bID2})
		require.Nil(t, err)

		keys, err := test.ust.GetKeys(username)
		require.Nil(t, err)
		require.NotNil(t, keys)
		require.Len(t, keys, 4)
	})
}
