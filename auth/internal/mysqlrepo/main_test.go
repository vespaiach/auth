package mysqlrepo

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/appconfig"
)

var repos *MysqlAppRepo

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	config := appconfig.LoadAppConfig()

	db, _ := initDb(config.DbConfig)
	defer db.Close()

	repos = NewMysqlAppRepo(db)

	code := m.Run()
	os.Exit(code)
}

func initDb(config *appconfig.DbConfig) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config.BuildMysqlDSN())

	if err != nil {
		return nil, err
	}

	return db, nil
}
