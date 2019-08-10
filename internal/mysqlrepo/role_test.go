package mysqlrepo

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	"github.com/vespaiach/auth/internal/comtype"
// )

// func loadRolefixtures() {
// 	tx := testApp.db.Begin()

// 	testApp.db.Exec("INSERT INTO `authentication`.`roles`(role_name, role_desc) VALUES (?, ?)", "admin", "admin")
// 	testApp.db.Exec("INSERT INTO `authentication`.`roles`(role_name, role_desc) VALUES (?, ?)", "staff", "staff")
// 	testApp.db.Exec("INSERT INTO `authentication`.`roles`(role_name, role_desc) VALUES (?, ?)", "guess", "guess")
// 	testApp.db.Exec("INSERT INTO `authentication`.`roles`(role_name, role_desc) VALUES (?, ?)", "admin_1", "admin 1")

// 	testApp.db.Exec("INSERT INTO `authentication`.`actions`(action_name, action_desc) VALUES (?, ?)", "create_user", "create_user")
// 	testApp.db.Exec("INSERT INTO `authentication`.`actions`(action_name, action_desc) VALUES (?, ?)", "update_user", "update_user")
// 	testApp.db.Exec("INSERT INTO `authentication`.`actions`(action_name, action_desc) VALUES (?, ?)", "delete_user", "delete_user")
// 	testApp.db.Exec("INSERT INTO `authentication`.`actions`(action_name, action_desc) VALUES (?, ?)", "query_user", "query_user")

// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 1, 1)
// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 1, 2)
// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 1, 3)
// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 1, 4)

// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 4, 1)
// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 4, 2)
// 	testApp.db.Exec("INSERT INTO `authentication`.`role_actions`(role_id, action_id) VALUES (?, ?)", 4, 3)

// 	tx.Commit()
// }

// func TestRoleRepo(tInst *testing.T) {

// 	s := &schema{
// 		up: `
// 			CREATE TABLE IF NOT EXISTS "authentication"."roles" (
// 				"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
// 				"role_name" VARCHAR(63) NOT NULL,
// 				"role_desc" VARCHAR(255) NOT NULL DEFAULT '',
// 				"active" TINYINT(1) NOT NULL DEFAULT 1,
// 				"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 				"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 			PRIMARY KEY ("id"),
// 			UNIQUE INDEX "role_name_uniq" ("role_name" ASC),
// 			INDEX "role_active_idx" ("active" ASC))
// 			ENGINE = InnoDB
// 			AUTO_INCREMENT = 1
// 			DEFAULT CHARACTER SET = utf8;

// 			CREATE TABLE IF NOT EXISTS "authentication"."actions" (
// 				"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
// 				"action_name" VARCHAR(63) NOT NULL,
// 				"action_desc" VARCHAR(255) NOT NULL DEFAULT '',
// 				"active" TINYINT(1) NOT NULL DEFAULT 1,
// 				"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 				"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
// 			PRIMARY KEY ("id"),
// 			UNIQUE INDEX "action_name_uniq" ("action_name" ASC),
// 			INDEX "action_active_idx" ("active" ASC))
// 			ENGINE = InnoDB
// 			AUTO_INCREMENT = 1
// 			DEFAULT CHARACTER SET = utf8;

// 			CREATE TABLE IF NOT EXISTS "authentication"."role_actions" (
// 				"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
// 				"role_id" BIGINT(20) UNSIGNED NOT NULL,
// 				"action_id" BIGINT(20) UNSIGNED NOT NULL,
// 				"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 			PRIMARY KEY ("id"),
// 			INDEX "action_id_idx" ("action_id" ASC),
// 			INDEX "role_id_idx" ("role_id" ASC),
// 			UNIQUE INDEX "role_action_uniq" ("role_id" ASC, "action_id" ASC),
// 			CONSTRAINT "action_id_on_permissions"
// 				FOREIGN KEY ("action_id")
// 				REFERENCES "authentication"."actions" ("id")
// 				ON DELETE CASCADE
// 				ON UPDATE CASCADE,
// 		    CONSTRAINT "role_id_on_permissions"
// 				FOREIGN KEY ("role_id")
// 				REFERENCES "authentication"."roles" ("id")
// 				ON DELETE CASCADE
// 				ON UPDATE CASCADE)
// 			ENGINE = InnoDB
// 			AUTO_INCREMENT = 1
// 			DEFAULT CHARACTER SET = utf8;
// 		`,
// 		down: `
// 			DROP TABLE IF EXISTS role_actions;
// 			DROP TABLE IF EXISTS roles;
// 			DROP TABLE IF EXISTS actions;
// 		`,
// 	}

// 	runWithSchema(s, tInst, func(t *testing.T) {

// 		loadRolefixtures()

// 		t.Run("query_role_by_name_success", func(t *testing.T) {
// 			filter := map[string]interface{}{"role_name": "admin"}
// 			roles, total, err := testApp.repos.RoleRepo.Query(1, 10, filter, map[string]comtype.SortDirection{})

// 			require.Nil(t, err)
// 			require.NotNil(t, roles)
// 			require.Len(t, roles, 2)
// 			require.Equal(t, int64(2), total)
// 			require.Greater(t, len(roles[0].Actions), 0)
// 			require.Greater(t, len(roles[1].Actions), 0)
// 		})

// 		t.Run("query_role_by_active_success", func(t *testing.T) {
// 			filter := map[string]interface{}{"active": comtype.Active}
// 			roles, total, err := testApp.repos.RoleRepo.Query(1, 10, filter, map[string]comtype.SortDirection{})

// 			require.Nil(t, err)
// 			require.NotNil(t, roles)
// 			require.Len(t, roles, 4)
// 			require.Equal(t, int64(4), total)
// 		})

// 		t.Run("query_role_by_active_fail", func(t *testing.T) {
// 			filter := map[string]interface{}{"active": comtype.Unactive}
// 			roles, total, err := testApp.repos.RoleRepo.Query(1, 10, filter, map[string]comtype.SortDirection{})

// 			require.Nil(t, err)
// 			require.NotNil(t, roles)
// 			require.Len(t, roles, 0)
// 			require.Equal(t, int64(0), total)
// 		})

// 		t.Run("create_role_success", func(t *testing.T) {
// 			id, err := testApp.repos.RoleRepo.Create("role_123", "Role 123")

// 			require.Nil(t, err)
// 			require.NotZero(t, id)
// 		})

// 		t.Run("get_role_id_success", func(t *testing.T) {
// 			role, err := testApp.repos.RoleRepo.GetByID(1)

// 			require.Nil(t, err)
// 			require.NotNil(t, role)
// 			require.Equal(t, "admin", role.RoleName)
// 			require.Equal(t, int64(1), role.ID)
// 		})

// 		t.Run("get_role_by_name_success", func(t *testing.T) {
// 			role, err := testApp.repos.RoleRepo.GetByName("admin_1")

// 			require.Nil(t, err)
// 			require.NotNil(t, role)
// 			require.Equal(t, role.RoleName, "admin_1")
// 			require.Equal(t, role.RoleDesc, "admin 1")
// 			require.Equal(t, role.Active, true)
// 			require.Len(t, role.Actions, 3)
// 		})

// 		t.Run("update_user_success", func(t *testing.T) {
// 			updating := map[string]interface{}{
// 				"role_name": "changed",
// 				"role_desc": "changed",
// 				"active":    false,
// 			}

// 			err := testApp.repos.RoleRepo.Update(int64(1), updating)
// 			require.Nil(t, err)
// 		})
// 	})
// }
