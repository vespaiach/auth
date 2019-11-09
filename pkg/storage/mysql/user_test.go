package mysql

import (
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/adding"
	"testing"
)

func TestStorage_AddUser(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_user", func(t *testing.T) {
		t.Parallel()

		username := test.mig.CreateUniqueString("username")
		email := test.mig.CreateUniqueString("email")

		id, err := test.st.AddUser(adding.User{username, email, "hash"})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_username", func(t *testing.T) {
		t.Parallel()

		username := test.mig.CreateUniqueString("username")
		email := test.mig.CreateUniqueString("email")

		test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["username"] = username
		})

		dupID, errDup := test.st.AddUser(adding.User{username, email, "hash"})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})

	t.Run("fail_add_a_duplicated_email", func(t *testing.T) {
		t.Parallel()

		username := test.mig.CreateUniqueString("username")
		email := test.mig.CreateUniqueString("email")

		test.mig.createSeedingUser(func(fields map[string]interface{}) {
			fields["email"] = email
		})

		dupID, errDup := test.st.AddUser(adding.User{username, email, "hash"})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}
