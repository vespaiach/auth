package mysql

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/pkg/common"
	"testing"
)

func TestBunchStorage_AddBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_add_a_bunch", func(t *testing.T) {
		t.Parallel()

		bunch := test.mig.createUniqueString("bunch")
		desc := test.mig.createUniqueString("desc")

		id, err := test.bst.AddBunch(bunch, desc)
		require.Nil(t, err)
		require.NotZero(t, id)
	})

	t.Run("fail_add_a_duplicated_bunch", func(t *testing.T) {
		t.Parallel()

		bunch := test.mig.createUniqueString("bunch")
		desc := test.mig.createUniqueString("desc")

		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = bunch
			fields["desc"] = desc
		})

		id, err := test.bst.AddBunch(bunch, desc)
		require.NotNil(t, err)
		require.Zero(t, id)
	})
}

func TestBunchStorage_GetBunchByName(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_bunch_by_name", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("bunch")

		id := test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name
		})

		bunch, err := test.bst.GetBunchByName(name)
		require.Nil(t, err)
		require.NotNil(t, bunch)
		require.Equal(t, id, bunch.ID)
	})
}

func TestBunchStorage_GetBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_get_a_bunch", func(t *testing.T) {
		t.Parallel()

		name := test.mig.createUniqueString("bunch")

		id := test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name
		})

		bunch, err := test.bst.GetBunch(id)
		require.Nil(t, err)
		require.NotNil(t, bunch)
		require.Equal(t, name, bunch.Name)
	})
}

func TestBunchStorage_ModifyBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_modify_a_bunch", func(t *testing.T) {
		t.Parallel()

		id := test.mig.createSeedingBunch(nil)
		name := test.mig.createUniqueString("name")
		desc := test.mig.createUniqueString("desc")

		err := test.bst.ModifyBunch(id, name, desc, sql.NullBool{
			Valid: true,
			Bool:  false,
		})
		require.Nil(t, err)

		bname, bdesc, bactive := test.mig.getBunchByID(id)
		require.Equal(t, name, bname)
		require.Equal(t, desc, bdesc)
		require.False(t, bactive)
	})
}

func TestBunchStorage_QueryBunches(t *testing.T) {
	t.Parallel()

	t.Run("success_modify_a_bunch", func(t *testing.T) {
		t.Parallel()

		name1 := test.mig.createUniqueString("tname")
		name2 := test.mig.createUniqueString("tname")
		name3 := test.mig.createUniqueString("tname")
		name4 := test.mig.createUniqueString("tname")
		name5 := test.mig.createUniqueString("tname")
		name6 := test.mig.createUniqueString("tname")

		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name1 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name2 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name3 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name4 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) { fields["name"] = name5 })
		test.mig.createSeedingBunch(func(fields map[string]interface{}) {
			fields["name"] = name6
			fields["active"] = false
		})

		rows, total, err := test.bst.QueryBunches(2, 2, "tname", sql.NullBool{Valid: false}, "created_at", common.Ascending)
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.Equal(t, int64(6), total)
		require.Len(t, rows, 2)

		rows, total, err = test.bst.QueryBunches(2, 2, "tname", sql.NullBool{Valid: true, Bool: true}, "created_at", common.Ascending)
		require.Nil(t, err)
		require.NotNil(t, rows)
		require.Equal(t, int64(5), total)
		require.Len(t, rows, 2)
	})
}

func TestBunchStorage_AddKeysToBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_keys_to_a_bunch", func(t *testing.T) {
		t.Parallel()

		kid1 := test.mig.createSeedingServiceKey(nil)
		kid2 := test.mig.createSeedingServiceKey(nil)
		kid3 := test.mig.createSeedingServiceKey(nil)
		kid4 := test.mig.createSeedingServiceKey(nil)
		bid := test.mig.createSeedingBunch(nil)

		err := test.bst.AddKeysToBunch(bid, []int64{kid1, kid2, kid3, kid4})
		require.Nil(t, err)

		results := test.mig.getKeyIDByBunchID(bid)
		require.NotNil(t, results)
		require.Len(t, results, 4)
	})
}

func TestBunchStorage_GetKeyInBunch(t *testing.T) {
	t.Parallel()

	t.Run("success_get_keys_in_a_bunch", func(t *testing.T) {
		t.Parallel()

		kid1 := test.mig.createSeedingServiceKey(nil)
		kid2 := test.mig.createSeedingServiceKey(nil)
		kid3 := test.mig.createSeedingServiceKey(nil)
		kid4 := test.mig.createSeedingServiceKey(nil)
		bid := test.mig.createSeedingBunch(nil)

		err := test.bst.AddKeysToBunch(bid, []int64{kid1, kid2, kid3, kid4})
		require.Nil(t, err)

		keys, err := test.bst.GetKeyInBunch(bid)
		require.Nil(t, err)
		require.NotNil(t, keys)
		require.Len(t, keys, 4)
	})
}

func TestBunchStorage_GetKeyIDs(t *testing.T) {
	t.Parallel()

	t.Run("success_get_keys_in_a_bunch", func(t *testing.T) {
		t.Parallel()

		name1 := test.mig.createUniqueString("tname")
		name2 := test.mig.createUniqueString("tname")
		name3 := test.mig.createUniqueString("tname")
		name4 := test.mig.createUniqueString("tname")
		name5 := test.mig.createUniqueString("tname")

		id1 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["key"] = name1 })
		id2 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["key"] = name2 })
		id3 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["key"] = name3 })
		id4 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["key"] = name4 })
		id5 := test.mig.createSeedingServiceKey(func(fields map[string]interface{}) { fields["key"] = name5 })

		ids, err := test.bst.GetKeyIDs([]string{name1, name2, name3, name4, name5})
		require.Nil(t, err)
		require.NotNil(t, ids)
		require.Len(t, ids, 5)
		require.True(t, contains(ids, id1))
		require.True(t, contains(ids, id2))
		require.True(t, contains(ids, id3))
		require.True(t, contains(ids, id4))
		require.True(t, contains(ids, id5))
	})
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
