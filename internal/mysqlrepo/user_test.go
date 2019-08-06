package mysqlrepo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vespaiach/auth/internal/comtype"
)

func loadUserfixtures() {
	tx := testApp.db.Begin()

	i := 0
	for i < 20 {
		a := string(i)
		testApp.db.Exec("INSERT INTO `authentication`.`users`(full_name, username, hashed, email) VALUES (?, ?, ?, ?)", "Toan "+a, "toan_"+a, "hased_"+a, "email_"+a)
		i++
	}

	tx.Commit()
}

func TestUserRepo(tInst *testing.T) {

	s := &schema{
		up: `
			CREATE TABLE IF NOT EXISTS "authentication"."users" (
				"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
				"full_name" VARCHAR(63) NOT NULL,
				"username" VARCHAR(63) NOT NULL,
				"hashed" VARCHAR(255) NOT NULL,
				"email" VARCHAR(127) NOT NULL,
				"active" TINYINT(1) NOT NULL DEFAULT 1,
				"verified" TINYINT(1) NOT NULL DEFAULT 0,
				"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY ("id"),
			UNIQUE INDEX "user_name_uniq" ("username" ASC),
			INDEX "active_idx" ("active" ASC))
			ENGINE = InnoDB
			AUTO_INCREMENT = 1
			DEFAULT CHARACTER SET = utf8;
		`,
		down: `
			DROP TABLE users;
		`,
	}

	runWithSchema(s, tInst, func(t *testing.T) {

		loadUserfixtures()

		t.Run("list_user_sucess", func(t *testing.T) {
			users, err := testApp.repos.UserRepo.Query(1, 10, map[string]interface{}{}, map[string]comtype.SortDirection{})

			require.Nil(t, err)
			require.NotNil(t, users)
			require.Equal(t, len(users), 10)
		})

		t.Run("create_user_sucess", func(t *testing.T) {
			user, err := testApp.repos.UserRepo.Create("Toan Nguyen", "username1", "123", "nta.toan@gmil.com")

			require.Nil(t, err)
			require.NotNil(t, user)
			require.NotZero(t, user.ID)
		})

	})
}
