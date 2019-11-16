package mysql

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage_AddServiceKey(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		id, err := test.kst.AddKey(key, desc)
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.createUniqueString("key")
		desc1 := test.mig.createUniqueString("desc")
		desc2 := test.mig.createUniqueString("desc")

		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["key"] = key
			fields["desc"] = desc1
		})

		dupID, errDup := test.kst.AddKey(key, desc2)
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}

func TestStorage_ModifyServiceKey(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_key", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingServiceKey(nil)
		key := test.mig.createUniqueString("key")
		desc := test.mig.createUniqueString("desc")

		err := test.kst.ModifyKey(id, key, desc)
		require.Nil(t, err)

		updatedKey, updatedDesc := test.mig.getServiceKeyByID(id)
		require.Equal(t, updatedDesc, desc)
		require.Equal(t, updatedKey, key)
	})
}