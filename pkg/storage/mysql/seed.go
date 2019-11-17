package mysql

var initDatabase = `
CREATE TABLE IF NOT EXISTS "keys" (
  	"id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	"key" VARCHAR(32) NOT NULL DEFAULT '',
	"desc" VARCHAR(64) NOT NULL DEFAULT '',
  	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  	"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "keys_key_uniq" ("key" ASC))
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "bunches" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "name" VARCHAR(32) NOT NULL DEFAULT '',
  "desc" VARCHAR(64) NOT NULL DEFAULT '',
  "active" TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "bunch_name_uniq" ("name" ASC),
  INDEX "bunch_active_idx" ("active" ASC))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS "users" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "username" VARCHAR(32) NOT NULL,
  "email" VARCHAR(64) NOT NULL,
  "hash" VARCHAR(255) NOT NULL,
  "active" TINYINT(1) NOT NULL DEFAULT 1,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "users_username_uniq" ("username" ASC),
  UNIQUE INDEX "users_email_uniq" ("email" ASC),
  INDEX "users_active_idx" ("active" ASC))
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "token_histories" (
  "uid" VARCHAR(36) NOT NULL,
  "user_id" BIGINT(20) UNSIGNED NOT NULL,
  "access_token" VARCHAR(1024) NOT NULL,
  "refresh_token" VARCHAR(1024) NOT NULL DEFAULT '',
  "remote_addr" VARCHAR(512) NOT NULL DEFAULT '',
  "x_forwarded_for" VARCHAR(512) NOT NULL DEFAULT '',
  "x_real_ip" VARCHAR(512) NOT NULL DEFAULT '',
  "user_agent" VARCHAR(512) NOT NULL DEFAULT '',
  "created_at" TIMESTAMP NOT NULL,
  "expired_at" TIMESTAMP NOT NULL,
  PRIMARY KEY ("uid"),
  UNIQUE INDEX "uid_uniq" ("uid" ASC))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "bunch_keys" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "bunch_id" BIGINT(20) UNSIGNED NOT NULL,
  "key_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "bunch_key_key_id_idx" ("key_id" ASC),
  INDEX "bunch_key_bunch_id_idx" ("bunch_id" ASC),
  UNIQUE INDEX "bunch_key_uniq" ("bunch_id" ASC, "key_id" ASC),
  CONSTRAINT "key_id_on_bunch_key"
    FOREIGN KEY ("key_id")
    REFERENCES "keys" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "role_id_on_bunch_key"
    FOREIGN KEY ("bunch_id")
    REFERENCES "bunches" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "user_bunches" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "user_id" BIGINT(20) UNSIGNED NOT NULL,
  "bunch_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "user_bunch_user_id_idx" ("user_id" ASC),
  INDEX "user_bunch_bunch_id_idx" ("bunch_id" ASC),
  UNIQUE INDEX "user_bunch_uniq" ("user_id" ASC, "bunch_id" ASC),
  CONSTRAINT "user_id_on_user_bunch"
    FOREIGN KEY ("user_id")
    REFERENCES "users" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "bunch_id_on_user_bunch"
    FOREIGN KEY ("bunch_id")
    REFERENCES "bunches" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB;
`

var dropDatabase = `
DROP TABLE IF EXISTS "user_bunches";
DROP TABLE IF EXISTS "bunch_keys";
DROP TABLE IF EXISTS "keys";
DROP TABLE IF EXISTS "bunches";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "token_histories";
`

var seedingData = `
INSERT INTO "keys" (id, "key", "desc") VALUES (1, 'create_action', 'Create a action');
INSERT INTO "keys" (id, "key", "desc") VALUES (2, 'delete_action', 'Delete a action');
INSERT INTO "keys" (id, "key", "desc") VALUES (3, 'update_action', 'Update a action');
INSERT INTO "keys" (id, "key", "desc") VALUES (4, 'view_action', 'View a action');
INSERT INTO "keys" (id, "key", "desc") VALUES (5 ,'list_action', 'List actions');
INSERT INTO "keys" (id, "key", "desc") VALUES (6, 'create_user', 'Create a user');
INSERT INTO "keys" (id, "key", "desc") VALUES (7, 'delete_user', 'Delete a user');
INSERT INTO "keys" (id, "key", "desc") VALUES (8, 'update_user', 'Update a user');
INSERT INTO "keys" (id, "key", "desc") VALUES (9, 'view_user', 'View a user');
INSERT INTO "keys" (id, "key", "desc") VALUES (10 ,'list_user', 'List users');
INSERT INTO bunches (id, "name", "desc") VALUES (1, 'admin_role', 'Admin role');
INSERT INTO bunches (id, "name", "desc") VALUES (2, 'staff_role', 'Staff role');
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (1, 1, 1);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (2, 1, 2);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (3, 1, 3);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (4, 1, 4);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (5, 1, 5);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (6, 1, 6);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (7, 1, 7);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (8, 1, 8);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (9, 1, 9);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (10, 1, 10);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (11, 2, 1);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (12, 2, 2);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (13, 2, 3);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (14, 2, 4);
INSERT INTO bunch_keys (id, bunch_id, key_id) VALUES (15, 2, 5);
INSERT INTO "users" (id, username, hash, email) VALUES (1, 'admin', '$2a$10$88y3eBfBma0lQzgEhPg7m.xmZQUE5DhcHqewtz0UvIYIfFZQFnD/G', 'admin@test.com');
INSERT INTO "users" (id, username, hash, email) VALUES (2, 'staff', '$2a$10$88y3eBfBma0lQzgEhPg7m.xmZQUE5DhcHqewtz0UvIYIfFZQFnD/G', 'staff@test.com');
INSERT INTO user_bunches (id, user_id, bunch_id) VALUES (1, 1, 1);
INSERT INTO user_bunches (id, user_id, bunch_id) VALUES (2, 1, 2);
INSERT INTO user_bunches (id, user_id, bunch_id) VALUES (3, 2, 2);
`
