package mysqlrepo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryTokenHistory(t *testing.T) {
	t.Parallel()

	t.Run("query_token_history_by_user_id_success", func(t *testing.T) {
		t.Parallel()

		histories, err := testApp.tokenHistoryRepo.GetByUserID(int64(1))

		require.Nil(t, err)
		require.NotNil(t, histories)
		require.Greater(t, len(histories), 0)
	})

	t.Run("query_token_history_by_user_id_empty", func(t *testing.T) {
		t.Parallel()

		histories, err := testApp.tokenHistoryRepo.GetByUserID(int64(0))

		require.Nil(t, err)
		require.NotNil(t, histories)
		require.Equal(t, len(histories), 0)
	})
}

func TestSaveTokenHistory(t *testing.T) {
	t.Parallel()

	t.Run("save_token_history_success", func(t *testing.T) {
		t.Parallel()

		err := testApp.tokenHistoryRepo.Save("e45e4c5a-c181-11e9-bf92-0242ac120002", 1, "access_token5", "refresh_token5",
			"localhost", "xForwardedFor", "xRealIP", "userAgent", time.Now(), time.Now().Add(time.Hour*2))
		require.Nil(t, err)
	})
}
