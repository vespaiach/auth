package mysqlrepo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

func TestQueryTokenHistory(t *testing.T) {
	t.Parallel()

	t.Run("query_token_history_by_user_id_success", func(t *testing.T) {
		t.Parallel()

		histories, total, err := testApp.tokenHistoryRepo.Query(1, 3, map[string]interface{}{
			"user_id": 1,
		}, map[string]comtype.SortDirection{
			"created_at": comtype.Ascending,
		})

		require.Nil(t, err)
		require.NotNil(t, histories)
		require.Greater(t, len(histories), 0)
		require.NotZero(t, total)
	})

	t.Run("query_token_history_by_date_success", func(t *testing.T) {
		t.Parallel()

		oneHourForward := time.Now().Add(time.Hour)
		twoHourForward := time.Now().Add(time.Hour * 2)
		threeHourForward := time.Now().Add(time.Hour * 3)

		testApp.createTokenHistory(func(h *model.TokenHistory) {
			h.CreatedAt = oneHourForward
			h.UserID = 2
		})
		testApp.createTokenHistory(func(h *model.TokenHistory) {
			h.CreatedAt = oneHourForward
			h.UserID = 2
		})
		testApp.createTokenHistory(func(h *model.TokenHistory) {
			h.CreatedAt = twoHourForward
			h.UserID = 2
		})
		testApp.createTokenHistory(func(h *model.TokenHistory) {
			h.CreatedAt = threeHourForward
			h.UserID = 2
		})

		histories, total, err := testApp.tokenHistoryRepo.Query(1, 4, map[string]interface{}{
			"from_date": oneHourForward.Add(-time.Second),
			"to_date":   twoHourForward.Add(time.Second),
			"user_id":   2,
		}, map[string]comtype.SortDirection{})

		fmt.Println(histories[0].CreatedAt)
		require.Nil(t, err)
		require.NotNil(t, histories)
		require.Len(t, histories, 3)
		require.Equal(t, int64(3), total)
	})
}

func TestCreateTokenHistory(t *testing.T) {
	t.Parallel()

	t.Run("create_token_history_success", func(t *testing.T) {
		t.Parallel()

		err := testApp.tokenHistoryRepo.Create("e45e4c5a-c181-11e9-bf92-0242ac120002", 1, "access_token5", "refresh_token5", time.Now())
		require.Nil(t, err)
	})
}
