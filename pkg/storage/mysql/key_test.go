package mysql

import (
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/adding"
	"testing"
)

func TestStorage_AddServiceKey(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.CreateUniqueString("key")
		desc := test.mig.CreateUniqueString("desc")

		id, err := test.st.AddServiceKey(adding.ServiceKey{key, desc})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_key", func(t *testing.T) {
		t.Parallel()

		key := test.mig.CreateUniqueString("key")
		desc1 := test.mig.CreateUniqueString("desc")
		desc2 := test.mig.CreateUniqueString("desc")

		test.mig.createSeedingServiceKey(func(fields map[string]interface{}) {
			fields["key"] = key
			fields["desc"] = desc1
		})

		dupID, errDup := test.st.AddServiceKey(adding.ServiceKey{key, desc2})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}
