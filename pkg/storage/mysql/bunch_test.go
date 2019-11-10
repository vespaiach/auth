package mysql

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/modifying"
	"testing"
)

func TestStorage_AddBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_bunch", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("name")
		desc := test.mig.createUniqueString("desc")

		id, err := test.st.AddBunch(adding.Bunch{name, desc})
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_bunch_name", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("name")
		desc1 := test.mig.createUniqueString("desc")
		desc2 := test.mig.createUniqueString("desc")

		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name
			fields["desc"] = desc1
		})

		dupID, errDup := test.st.AddBunch(adding.Bunch{name, desc2})
		require.NotNil(t, errDup)
		require.Zero(t, dupID)
	})
}

func TestStorage_ModifyBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_update_a_bunch", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingBunch(nil)

		err := test.st.ModifyBunch(modifying.Bunch{
			ID:   id,
			Name: "updated",
			Desc: "updated",
			Active: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		})
		require.Nil(t, err)

		name, desc, active := test.mig.getBunchByID(id)
		require.Equal(t, name, "updated")
		require.Equal(t, desc, "updated")
		require.Equal(t, active, true)
	})
}
