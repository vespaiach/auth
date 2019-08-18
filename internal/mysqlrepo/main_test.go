package mysqlrepo

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/migrate"
	"github.com/vespaiach/auth/internal/model"
)

type appTesting struct {
	actionRepo       model.ActionRepo
	roleRepo         model.RoleRepo
	userRepo         model.UserRepo
	userActionRepo   model.UserActionRepo
	userRoleRepo     model.UserRoleRepo
	roleActionRepo   model.RoleActionRepo
	tokenHistoryRepo model.TokenHistoryRepo
	config           *conf.AppConfig
	actionIDs        []int64
	db               *sqlx.DB
}

var testApp *appTesting

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	config := conf.LoadAppConfig()

	db, _ := initDb(config.DbConfig)

	testApp = new(appTesting)
	testApp.actionRepo = NewMysqlActionRepo(db)
	testApp.roleRepo = NewMysqlRoleRepo(db)
	testApp.userRepo = NewMysqlUserRepo(db)
	testApp.userActionRepo = NewMysqlUserActionRepo(db)
	testApp.userRoleRepo = NewMysqlUserRoleRepo(db)
	testApp.roleActionRepo = NewMysqlRoleActionRepo(db)
	testApp.tokenHistoryRepo = NewMysqlTokenHistoryRepo(db)
	testApp.config = config
	testApp.db = db

	mig := migrate.NewMigrator(db)

	mig.Up()
	mig.SeedTestData()

	code := m.Run()

	mig.Down()
	testApp.db.Close()
	os.Exit(code)
}

func initDb(config *conf.DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.BuildMysqlDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
