package mysqlrepo

import (
	"strings"
	"testing"
	"time"

	"github.com/vespaiach/auth/internal/comtype"

	"github.com/stretchr/testify/require"
)

func TestSQLWhereBuilder(t *testing.T) {
	t.Parallel()

	t.Run("empty_condition_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlWhereBuilder(" AND ", map[string]interface{}{})

		require.Len(t, sql, 0)
	})

	t.Run("one_condition_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlWhereBuilder(" AND ", map[string]interface{}{
			"user_id": true,
		})

		require.Equal(t, "WHERE user_id = :user_id", sql)
	})

	t.Run("two_conditions_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlWhereBuilder(" AND ", map[string]interface{}{
			"user_id": 1,
			"name":    "toan",
		})

		require.True(t, strings.Contains(sql, ":user_id"))
		require.True(t, strings.Contains(sql, "LIKE :name"))
		require.Equal(t, strings.Count(sql, "AND"), 1)
	})

	t.Run("date_conditions_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlWhereBuilder(" AND ", map[string]interface{}{
			"to_date":   time.Now().Format(comtype.DateTimeLayout),
			"from_date": time.Now().Format(comtype.DateTimeLayout),
		})

		require.True(t, strings.Contains(sql, "created_at >= :from_date"))
		require.True(t, strings.Contains(sql, "created_at <= :to_date"))
	})
}

func TestSortingBuilder(t *testing.T) {
	t.Parallel()

	t.Run("empty_condition_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlSortingBuilder(map[string]comtype.SortDirection{})

		require.Equal(t, sql, "created_at DESC")
	})

	t.Run("one_field_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlSortingBuilder(map[string]comtype.SortDirection{
			"name": comtype.Decending,
		})

		require.Equal(t, "name DESC", sql)
	})

	t.Run("two_fields_create_success", func(t *testing.T) {
		t.Parallel()

		sql := sqlSortingBuilder(map[string]comtype.SortDirection{
			"user_id": comtype.Decending,
			"name":    comtype.Ascending,
		})

		require.True(t, strings.Contains(sql, "user_id DESC"))
		require.True(t, strings.Contains(sql, "name ASC"))
	})
}
