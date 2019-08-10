package mysqlrepo

// import (
// 	"strconv"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	"github.com/vespaiach/auth/internal/comtype"
// 	"github.com/vespaiach/auth/internal/model"
// )

// func loadUserfixtures() {
// 	tx := testApp.db.Begin()

// 	i := 0
// 	for i < 20 {
// 		a := strconv.Itoa(i)
// 		testApp.db.Exec("INSERT INTO `authentication`.`users`(full_name, username, hashed, email) VALUES (?, ?, ?, ?)", "Toan "+a, "toan_"+a, "hashed_"+a, "email_"+a)
// 		i++
// 	}

// 	tx.Commit()
// }

// func TestUserRepo(tInst *testing.T) {

// 	s := &schema{
// 		up: `
// 			CREATE TABLE IF NOT EXISTS "authentication"."users" (
// 				"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
// 				"full_name" VARCHAR(63) NOT NULL,
// 				"username" VARCHAR(63) NOT NULL,
// 				"hashed" VARCHAR(255) NOT NULL,
// 				"email" VARCHAR(127) NOT NULL,
// 				"active" TINYINT(1) NOT NULL DEFAULT 1,
// 				"verified" TINYINT(1) NOT NULL DEFAULT 0,
// 				"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 				"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 			PRIMARY KEY ("id"),
// 			UNIQUE INDEX "user_name_uniq" ("username" ASC),
// 			INDEX "active_idx" ("active" ASC))
// 			ENGINE = InnoDB
// 			AUTO_INCREMENT = 1
// 			DEFAULT CHARACTER SET = utf8;
// 		`,
// 		down: `
// 			DROP TABLE users;
// 		`,
// 	}

// 	runWithSchema(s, tInst, func(t *testing.T) {

// 		loadUserfixtures()

// 		var user1 *model.User

// 		t.Run("query_user_by_full_name_success", func(t *testing.T) {
// 			filter := map[string]interface{}{"full_name": "Toan 1"}
// 			users, total, err := testApp.repos.UserRepo.Query(1, 10, filter, map[string]comtype.SortDirection{})

// 			require.Nil(t, err)
// 			require.NotNil(t, users)
// 			require.Len(t, users, 10)
// 			require.Equal(t, total, int64(11))
// 		})

// 		t.Run("query_user_by_email_success", func(t *testing.T) {
// 			filter := map[string]interface{}{"email": "Ema"}
// 			users, total, err := testApp.repos.UserRepo.Query(1, 5, filter, map[string]comtype.SortDirection{})

// 			require.Nil(t, err)
// 			require.NotNil(t, users)
// 			require.Len(t, users, 5)
// 			require.Equal(t, total, int64(20))
// 		})

// 		t.Run("create_user_success", func(t *testing.T) {
// 			user, err := testApp.repos.UserRepo.Create("Toan Nguyen", "username1", "123", "nta.toan@gmil.com")

// 			require.Nil(t, err)
// 			require.NotNil(t, user)
// 			require.NotZero(t, user.ID)
// 		})

// 		t.Run("get_user_by_username_success", func(t *testing.T) {
// 			var er error
// 			user1, er = testApp.repos.UserRepo.GetByUsername("toan_1")

// 			require.Nil(t, er)
// 			require.NotNil(t, user1)
// 			require.Equal(t, user1.FullName, "Toan 1")
// 			require.Equal(t, user1.Hashed, "hashed_1")
// 			require.Equal(t, user1.Email, "email_1")
// 		})

// 		t.Run("get_user_by_id_success", func(t *testing.T) {
// 			user, err := testApp.repos.UserRepo.GetByID(user1.ID)

// 			require.Nil(t, err)
// 			require.NotNil(t, user)
// 			require.Equal(t, user.ID, user1.ID)
// 		})

// 		t.Run("get_user_by_email_success", func(t *testing.T) {
// 			user, err := testApp.repos.UserRepo.GetByEmail(user1.Email)

// 			require.Nil(t, err)
// 			require.NotNil(t, user)
// 			require.Equal(t, user.ID, user1.ID)
// 		})

// 		t.Run("update_user_success", func(t *testing.T) {
// 			updating := map[string]interface{}{
// 				"full_name": "changed",
// 				"email":     "changed",
// 				"hashed":    "changed",
// 				"active":    false,
// 				"verified":  true,
// 			}
// 			err := testApp.repos.UserRepo.Update(user1.ID, updating)

// 			require.Nil(t, err)
// 		})
// 	})
// }
