package migrate

var initDatabase = `
CREATE TABLE IF NOT EXISTS "actions" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
	"action_name" VARCHAR(63) NOT NULL DEFAULT '',
	"action_desc" VARCHAR(254) NOT NULL DEFAULT '',
  "active" TINYINT(1) NOT NULL DEFAULT 1,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "actions_action_name_uniq" ("action_name" ASC),
  INDEX "actions_active_idx" ("active" ASC))
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "roles" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "role_name" VARCHAR(63) NOT NULL DEFAULT '',
  "role_desc" VARCHAR(254) NOT NULL DEFAULT '',
  "active" TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE INDEX "roles_role_name_uniq" ("role_name" ASC),
  INDEX "roles_active_idx" ("active" ASC))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS "users" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "full_name" VARCHAR(255) NOT NULL,
  "username" VARCHAR(63) NOT NULL,
  "hashed" VARCHAR(255) NOT NULL,
  "email" VARCHAR(127) NOT NULL,
  "verified" TINYINT(1) NOT NULL DEFAULT 0,
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
  "refresh_token" VARCHAR(1024) NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("uid"),
  UNIQUE INDEX "uid_uniq" ("uid" ASC))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "role_actions" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "role_id" BIGINT(20) UNSIGNED NOT NULL,
  "action_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "role_actions_action_id_idx" ("action_id" ASC),
  INDEX "role_actions_role_id_idx" ("role_id" ASC),
  UNIQUE INDEX "role_actions_uniq" ("role_id" ASC, "action_id" ASC),
  CONSTRAINT "action_id_on_role_actions"
    FOREIGN KEY ("action_id")
    REFERENCES "actions" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "role_id_on_role_actions"
    FOREIGN KEY ("role_id")
    REFERENCES "roles" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
AUTO_INCREMENT = 1
DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS "user_roles" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "user_id" BIGINT(20) UNSIGNED NOT NULL,
  "role_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "user_roles_user_id_idx" ("user_id" ASC),
  INDEX "user_roles_role_id_idx" ("role_id" ASC),
  UNIQUE INDEX "user_roles_uniq" ("user_id" ASC, "role_id" ASC),
  CONSTRAINT "user_id_on_user_roles"
    FOREIGN KEY ("user_id")
    REFERENCES "users" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "role_id_on_user_roles"
    FOREIGN KEY ("role_id")
    REFERENCES "roles" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS "user_actions" (
  "id" BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  "user_id" BIGINT(20) UNSIGNED NOT NULL,
  "action_id" BIGINT(20) UNSIGNED NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  INDEX "user_actions_user_id_idx" ("user_id" ASC),
  INDEX "user_actions_action_id_idx" ("action_id" ASC),
  UNIQUE INDEX "user_actions_uniq" ("user_id" ASC, "action_id" ASC),
  CONSTRAINT "user_id_on_user_actions"
    FOREIGN KEY ("user_id")
    REFERENCES "users" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT "action_id_on_user_actions"
    FOREIGN KEY ("action_id")
    REFERENCES "actions" ("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB;
`

var dropInitDatabase = `
DROP TABLE IF EXISTS "user_actions";
DROP TABLE IF EXISTS "role_actions";
DROP TABLE IF EXISTS "user_roles";
DROP TABLE IF EXISTS "actions";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "token_histories";
`
