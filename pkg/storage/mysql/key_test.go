package mysql

import (
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/storage"
	"testing"
)

func TestKeyMysqlStorer_Insert(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		id, err := test.kst.Insert(storage.Key{Name: key, Desc: desc})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc1 := test.mig.createUniqueString("desc")
		desc2 := test.mig.createUniqueString("desc")

		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["name"] = key
			fields["desc"] = desc1
		})

		dupID, errDup := test.kst.Insert(storage.Key{Name: key, Desc: desc2})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}

func TestKeyMysqlStorer_Update(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_key", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingServiceKey(nil)
		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		err := test.kst.Update(storage.Key{
			ID:   id,
			Name: key,
			Desc: desc,
		})
		require.Nil(t, err)
	})
}

func TestKeyMysqlStorer_Delete(t *testing.T) {
	t.Parallel()

	t.Run("success_delete_a_key", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingServiceKey(nil)

		err := test.kst.Delete(id)
		require.Nil(t, err)
	})
}

func TestKeyMysqlStorer_Get(t *testing.T) {
	t.Parallel()

	t.Run("success_delete_a_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		id := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["name"] = key
			fields["desc"] = desc
		})

		found, err := test.kst.Get(id)
		require.Nil(t, err)
		require.NotNil(t, found)
		require.Equal(t, found.ID, id)
		require.Equal(t, found.Name, key)
		require.Equal(t, found.Desc, desc)
	})
}
