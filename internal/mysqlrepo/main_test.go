package mysqlrepo

import (
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/conf"
)

type TestApp struct {
	repos  *MysqlAppRepo
	config *conf.AppConfig
	db     *gorm.DB
}

type schema struct {
	up   string
	down string
}

var testApp *TestApp

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	config := conf.LoadAppConfig()

	db, _ := initDb(config.DbConfig)
	defer db.Close()

	repos := NewMysqlAppRepo(db)

	testApp = new(TestApp)
	testApp.repos = repos
	testApp.config = config
	testApp.db = db

	code := m.Run()
	os.Exit(code)
}

func execSQL(query string) {
	stmts := strings.Split(query, ";\n")
	if len(strings.Trim(stmts[len(stmts)-1], " \n\t\r")) == 0 {
		stmts = stmts[:len(stmts)-1]
	}
	for _, s := range stmts {
		testApp.db.Exec(s)

		if len(testApp.db.GetErrors()) > 0 {
			fmt.Println("error in running schema creating")
		}
	}
}

func initDb(config *conf.DbConfig) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config.BuildMysqlDSN())
	db.LogMode(true)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s schema) santize() (string, string) {
	return strings.Replace(s.up, `"`, "`", -1), s.down
}

func runWithSchema(s *schema, t *testing.T, testGroup func(t *testing.T)) {
	up, down := s.santize()
	defer func() {
		execSQL(down)
	}()

	execSQL(up)

	testGroup(t)
}
