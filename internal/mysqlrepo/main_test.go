package mysqlrepo

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

type appTesting struct {
	actionRepo model.ActionRepo
	roleRepo   model.RoleRepo
	config     *conf.AppConfig
	actionIDs  []int64
	db         *sqlx.DB
}

var testApp *appTesting

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	config := conf.LoadAppConfig()

	db, _ := initDb(config.DbConfig)

	testApp = new(appTesting)
	testApp.actionRepo = NewMysqlActionRepo(db)
	testApp.roleRepo = NewMysqlRoleRepo(db)
	testApp.config = config
	testApp.db = db

	err := testApp.createActionTable()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = testApp.createRoleTable()
	if err != nil {
		log.Fatal(err)
		return
	}

	code := m.Run()

	cleanUp()
	os.Exit(code)
}

func initDb(config *conf.DbConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.BuildMysqlDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func cleanUp() {
	fmt.Println("-----------clean-up-----------")

	testApp.dropActionTable()
	testApp.dropRoleTable()
	testApp.db.Close()
}
