package mysql

import (
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/adding"
	"testing"
)

func TestStorage_AddBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_bunch", func(t *testing.T) {
		t.Parallel()

		name := test.mig.CreateUniqueString("name")
		desc := test.mig.CreateUniqueString("desc")

		id, err := test.st.AddBunch(adding.Bunch{name, desc})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_bunch_name", func(t *testing.T) {
		t.Parallel()

		name := test.mig.CreateUniqueString("name")
		desc1 := test.mig.CreateUniqueString("desc")
		desc2 := test.mig.CreateUniqueString("desc")

		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name
			fields["desc"] = desc1
		})

		dupID, errDup := test.st.AddBunch(adding.Bunch{name, desc2})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}
