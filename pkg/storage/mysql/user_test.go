package mysql

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/modifying"
	"testing"
)

func TestStorage_AddUser(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_user", func(t *testing.T) {
		t.Parallel()

		username := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")

		id, err := test.st.AddUser(adding.User{username, email, "hash"})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_username", func(t *testing.T) {
		t.Parallel()

		username := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")

		test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["username"] = username
		})

		dupID, errDup := test.st.AddUser(adding.User{username, email, "hash"})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})

	t.Run("fail_add_a_duplicated_email", func(t *testing.T) {
		t.Parallel()

		username := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")

		test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["email"] = email
		})

		dupID, errDup := test.st.AddUser(adding.User{username, email, "hash"})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}

func TestStorage_ModifyUser(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_user", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingUser(nil)

		err := test.st.ModifyUser(modifying.User{
			ID:       id,
			Username: "updated",
			Email:    "updated",
			Active: sql.NullBool{
				false,
				true,
			},
		})
		require.Nil(t, err)

		username, email, active := test.mig.getUserByID(id)
		require.Equal(t, username, "updated")
		require.Equal(t, email, "updated")
		require.Equal(t, active, false)
	})
}

func TestStorage_GetUserByID(t *testing.T) {
	t.Parallel()

	t.Run("success_get_user_by_id", func(t *testing.T) {
		t.Parallel()

		username := test.mig.createUniqueString("username")
		email := test.mig.createUniqueString("email")

		id := test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["username"] = username
			fields["email"] = email
		})

		user, err := test.st.GetUserByID(id)
		require.Nil(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.ID, id)
		require.Equal(t, user.Username, username)
		require.Equal(t, user.Email, email)
		require.True(t, user.Active)
	})
}
