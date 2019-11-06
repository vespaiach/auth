package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func TestRegisterUser(t *testing.T) {
	t.Run("register_new_user_success", func(t *testing.T) {
		user, err := testApp.userService.RegisterUser("Toan Nguyen", "vespa", "hashed", "toan@test.com")

		require.Nil(t, err)
		require.NotNil(t, user)
		require.NotZero(t, user.ID)
	})

	t.Run("duplicated_user_name", func(t *testing.T) {
		email := testApp.mig.CreateUniqueString("email")
		user := testApp.mig.CreateSeedingUser(nil)

		dup, err := testApp.userService.RegisterUser("Toan Nguyen", user.Username, "hashed", email)
		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDuplicatedData))
		require.Nil(t, dup)
	})

	t.Run("duplicated_email", func(t *testing.T) {
		username := testApp.mig.CreateUniqueString("email")
		user := testApp.mig.CreateSeedingUser(nil)

		dup, err := testApp.userService.RegisterUser("Toan Nguyen", username, "password", user.Email)
		require.NotNil(t, err)
		require.True(t, err.Is(comtype.ErrDuplicatedData))
		require.Nil(t, dup)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("update_user_password", func(t *testing.T) {
		username := testApp.mig.CreateUniqueString("username")
		user := testApp.mig.CreateSeedingUser(func(fields map[string]interface{}) {
			fields["hashed"] = "pass"
			fields["username"] = username
		})

		updated, err := testApp.userService.UpdateUser(user.ID, "", "", "", "password", nil, nil)
		require.Nil(t, err)
		require.NotNil(t, updated)
		require.NotZero(t, updated.ID)

		verified, actions, roles, err := testApp.userService.VerifyLogin(username, "password")
		require.Nil(t, err)
		require.NotNil(t, verified)
		require.NotNil(t, actions)
		require.NotNil(t, roles)
		require.NotZero(t, verified.ID)
	})

	t.Run("update_other_user_data", func(t *testing.T) {
		user := testApp.mig.CreateSeedingUser(nil)
		status := false

		updated, err := testApp.userService.UpdateUser(user.ID, user.FullName+"_updated", user.Username+"_updated",
			user.Email+"_updated", "", &status, &status)
		require.Nil(t, err)
		require.NotNil(t, updated)
		require.Equal(t, updated.FullName, user.FullName+"_updated")
		require.Equal(t, updated.Username, user.Username+"_updated")
		require.Equal(t, updated.Email, user.Email+"_updated")
		require.False(t, updated.Verified)
		require.False(t, updated.Active)
	})

	t.Run("update_duplicated_email_username", func(t *testing.T) {
		user := testApp.mig.CreateSeedingUser(nil)

		dupUsername, err := testApp.userService.UpdateUser(user.ID, "", user.Username, "", "", nil, nil)
		require.NotNil(t, err)
		require.Nil(t, dupUsername)
		require.True(t, err.Is(comtype.ErrDuplicatedData))

		dupEmail, err := testApp.userService.UpdateUser(user.ID, "", "", user.Email, "", nil, nil)
		require.NotNil(t, err)
		require.Nil(t, dupEmail)
		require.True(t, err.Is(comtype.ErrDuplicatedData))
	})
}

func TestVerifyUser(t *testing.T) {

	t.Run("verify_login_success", func(t *testing.T) {
		username := testApp.mig.CreateUniqueString("username")
		testApp.mig.CreateSeedingUser(func(fields map[string]interface{}) {
			fields["username"] = username
		})

		found, actions, roles, err := testApp.userService.VerifyLogin(username, "password")
		require.Nil(t, err)
		require.NotNil(t, found)
		require.NotNil(t, actions)
		require.NotNil(t, roles)
		require.NotZero(t, found.ID)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("get_user_by_id", func(t *testing.T) {
		user := testApp.mig.CreateSeedingUser(nil)

		found, err := testApp.userService.GetUser(user.ID)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.Username, user.Username)
		require.Equal(t, found.FullName, user.FullName)
		require.Equal(t, found.Email, user.Email)
	})
}

func TestQueryUsers(t *testing.T) {
	t.Parallel()

	t.Run("query_users", func(t *testing.T) {
		testApp.mig.CreateSeedingUser(nil)
		testApp.mig.CreateSeedingUser(nil)
		testApp.mig.CreateSeedingUser(nil)

		users, err := testApp.userService.FetchUsers(3, "", "", nil, "", nil, "-full_name")
		require.Nil(t, err)
		require.NotNil(t, users)
		require.Len(t, users, 3)
	})
}
