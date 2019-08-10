package mysqlrepo

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func loadActionfixtures() {
	tx := testApp.db.Begin()

	i := 0
	for i < 20 {
		a := strconv.Itoa(i)
		testApp.db.Exec("INSERT INTO `authentication`.`actions`(action_name, action_desc) VALUES (?, ?)", "toan_name_"+a, "toan_desc_"+a)
		i++
	}

	tx.Commit()
}

func TestActionRepo(tInst *testing.T) {

	s := &schema{
		up: `
			CREATE TABLE IF NOT EXISTS "authentication"."actions" (
				"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
				"action_name" VARCHAR(63) NOT NULL,
				"action_desc" VARCHAR(255) NOT NULL DEFAULT '',
				"active" TINYINT(1) NOT NULL DEFAULT 1,
				"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY ("id"),
			UNIQUE INDEX "action_name_uniq" ("action_name" ASC),
			INDEX "action_active_idx" ("active" ASC))
			ENGINE = InnoDB
			AUTO_INCREMENT = 1
			DEFAULT CHARACTER SET = utf8;
		`,
		down: `
			DROP TABLE actions;
		`,
	}

	runWithSchema(s, tInst, func(t *testing.T) {

		loadActionfixtures()

		// var user1 *model.User

		t.Run("query_action_by_name_success", func(t *testing.T) {
			filter := map[string]interface{}{"action_name": "toan_name_"}
			actions, total, err := testApp.repos.ActionRepo.Query(1, 10, filter, map[string]comtype.SortDirection{})

			require.Nil(t, err)
			require.NotNil(t, actions)
			require.Len(t, actions, 10)
			require.Equal(t, total, int64(20))
		})

		t.Run("query_action_by_active_success", func(t *testing.T) {
			filter := map[string]interface{}{"active": false}
			actions, total, err := testApp.repos.ActionRepo.Query(1, 10, filter, map[string]comtype.SortDirection{})

			require.Nil(t, err)
			require.NotNil(t, actions)
			require.Len(t, actions, 0)
			require.Equal(t, total, int64(0))
		})

		t.Run("create_action_success", func(t *testing.T) {
			action, err := testApp.repos.ActionRepo.Create("action_123", "Action 123")

			require.Nil(t, err)
			require.NotNil(t, action)
			require.NotZero(t, action.ID)
		})

		// t.Run("get_user_by_username_success", func(t *testing.T) {
		// 	var er error
		// 	user1, er = testApp.repos.ActionRepo.GetByUsername("toan_1")

		// 	require.Nil(t, er)
		// 	require.NotNil(t, user1)
		// 	require.Equal(t, user1.FullName, "Toan 1")
		// 	require.Equal(t, user1.Hashed, "hashed_1")
		// 	require.Equal(t, user1.Email, "email_1")
		// })

		// t.Run("get_user_by_id_success", func(t *testing.T) {
		// 	user, err := testApp.repos.ActionRepo.GetByID(user1.ID)

		// 	require.Nil(t, err)
		// 	require.NotNil(t, user)
		// 	require.Equal(t, user.ID, user1.ID)
		// })

		// t.Run("get_user_by_email_success", func(t *testing.T) {
		// 	user, err := testApp.repos.ActionRepo.GetByEmail(user1.Email)

		// 	require.Nil(t, err)
		// 	require.NotNil(t, user)
		// 	require.Equal(t, user.ID, user1.ID)
		// })

		// t.Run("update_user_success", func(t *testing.T) {
		// 	updating := map[string]interface{}{
		// 		"full_name": "changed",
		// 		"email":     "changed",
		// 		"hashed":    "changed",
		// 		"active":    false,
		// 		"verified":  true,
		// 	}
		// 	err := testApp.repos.ActionRepo.Update(user1.ID, updating)

		// 	require.Nil(t, err)
		// })
	})
}
