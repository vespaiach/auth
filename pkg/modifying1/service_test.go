package modifying

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_ModifyServiceKey(t *testing.T) {
	t.Parallel()

	t.Run("success_modify_a_key", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyServiceKey(ServiceKey{
			ID:   10,
			Key:  "test",
			Desc: "test",
		})
		require.Nil(t, err)

		err = testService.ModifyServiceKey(ServiceKey{
			ID:   10,
			Key:  "fake_key",
			Desc: "test",
		})
		require.Nil(t, err)

		err = testService.ModifyServiceKey(ServiceKey{
			ID:   10,
			Key:  "",
			Desc: "",
		})
		require.Nil(t, err)
	})

	t.Run("fail_to_modify_not_existing_id_key", func(t *testing.T) {
		t.Parallel()

		sk := ServiceKey{
			ID:   1,
			Key:  "test",
			Desc: "",
		}

		err := testService.ModifyServiceKey(sk)
		require.NotNil(t, err)
		require.Equal(t, err.Error(), "no key found with id = 1")
	})

	t.Run("fail_to_modify_duplicated_key", func(t *testing.T) {
		t.Parallel()

		sk := ServiceKey{
			ID:   10,
			Key:  "dup",
			Desc: "test",
		}

		err := testService.ModifyServiceKey(sk)
		require.NotNil(t, err)
		require.Equal(t, err, ErrDuplicatedKey)
	})

	t.Run("fail_to_modify_key_validate", func(t *testing.T) {
		t.Parallel()

		sk := ServiceKey{
			ID:   10,
			Key:  "dup ss fas fas asd asd asdfas asdf asf as asdfoiwuer weori qowir oqwi row rowi roiworqwodko",
			Desc: "",
		}

		err := testService.ModifyServiceKey(sk)
		require.NotNil(t, err)
		require.Equal(t, err, ErrServiceKeyTooLong)
	})
}

func TestService_ModifyBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_modify_a_bunch", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyBunch(Bunch{
			ID:     10,
			Name:   "",
			Desc:   "",
			Active: sql.NullBool{},
		})
		require.Nil(t, err)

		err = testService.ModifyBunch(Bunch{
			ID:   10,
			Name: "fake_bunch",
			Desc: "test",
		})
		require.Nil(t, err)

		err = testService.ModifyBunch(Bunch{
			ID:   10,
			Name: "",
			Desc: "dup",
		})
		require.Nil(t, err)
	})

	t.Run("fail_to_modify_not_existing_id", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyBunch(Bunch{
			ID:   1,
			Name: "test",
			Desc: "",
		})
		require.NotNil(t, err)
		require.Equal(t, err.Error(), "no bunch found with id = 1")
	})

	t.Run("fail_to_modify_duplicated_bunch", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyBunch(Bunch{
			ID:   10,
			Name: "dup",
			Desc: "test",
		})
		require.NotNil(t, err)
		require.Equal(t, err, ErrDuplicatedBunch)
	})

	t.Run("fail_to_modify_key_validate", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyBunch(Bunch{
			ID:   10,
			Name: "dup ss fas fas asd asd asdfas asdf asf as asdfoiwuer weori qowir oqwi row rowi roiworqwodko",
			Desc: "",
		})
		require.NotNil(t, err)
		require.Equal(t, err, ErrBunchNameTooLong)
	})
}

func TestService_ModifyUser(t *testing.T) {
	t.Parallel()

	t.Run("success_modify_a_user", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyUser(User{
			ID:       10,
			Username: "",
			Email:    "",
			Active:   sql.NullBool{},
		})
		require.Nil(t, err)

		err = testService.ModifyUser(User{
			ID:       10,
			Username: "fake_username",
			Email:    "",
			Active:   sql.NullBool{},
		})
		require.Nil(t, err)

		err = testService.ModifyUser(User{
			ID:       10,
			Username: "",
			Email:    "fake_email@gmail.com",
			Active:   sql.NullBool{},
		})
		require.Nil(t, err)
	})

	t.Run("fail_to_modify_not_existing_id", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyUser(User{
			ID:       1,
			Username: "test",
			Email:    "",
		})
		require.NotNil(t, err)
		require.Equal(t, err.Error(), "no user found with id = 1")
	})

	t.Run("fail_to_modify_duplicated_username", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyUser(User{
			ID:       10,
			Username: "dup",
			Email:    "",
		})
		require.NotNil(t, err)
		require.Equal(t, err, ErrDuplicatedUsername)
	})

	t.Run("fail_to_modify_duplicated_email", func(t *testing.T) {
		t.Parallel()

		err := testService.ModifyUser(User{
			ID:       10,
			Username: "",
			Email:    "dup",
		})
		require.NotNil(t, err)
		require.Equal(t, err, ErrDuplicatedEmail)
	})
}
