package service

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/migrate"
	"github.com/vespaiach/auth/internal/mysqlrepo"
)

type appTesting struct {
	userService       *UserService
	tokenService      TokenService
	actionService     *ActionService
	roleService       *RoleService
	roleActionService *RoleActionService
	userRoleService   *UserRoleService
	config            *conf.AppConfig
	actionIDs         []int64
	db                *sqlx.DB
	mig               *migrate.Migrator
}

var testApp *appTesting

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	config := conf.LoadAppConfig()

	db, _ := initDb(config)

	appRepo := mysqlrepo.NewMysqlAppRepo(db)
	testApp = new(appTesting)
	testApp.userService = NewUserService(appRepo, config)
	testApp.tokenService = NewTokenService(appRepo, config)
	testApp.actionService = NewActionService(appRepo, config)
	testApp.roleService = NewRoleService(appRepo, config)
	testApp.roleActionService = NewRoleActionService(appRepo, config)
	testApp.userRoleService = NewUserRoleService(appRepo, config)
	testApp.config = config
	testApp.db = db

	mig := migrate.NewMigrator(db)
	testApp.mig = mig

	mig.Down()
	mig.Up()
	mig.SeedTestData()

	code := m.Run()

	mig.Down()
	testApp.db.Close()
	os.Exit(code)
}

func initDb(config *conf.AppConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.BuildMysqlDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
