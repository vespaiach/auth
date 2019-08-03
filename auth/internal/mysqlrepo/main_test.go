package mysqlrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/jinzhu/gorm"
	"github.com/vespaiach/auth/internal/appconfig"
)

var repos *MysqlAppRepo
var config *appconfig.AppConfig

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	config = appconfig.LoadAppConfig()

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

func prepareFixtures(testName string) *sql.DB {
	db, err := sql.Open("mysql", config.DbConfig.BuildMysqlDSN())
	if err != nil {
		log.Fatalf("could not connect to the MySQL database... %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	// Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path.Join(config.CommonConfig.AppDir, "internal/mysqlrepo/fixtures", testName)), // file://path/to/directory
		"mysql", driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")

	return db
}

func clearFixtures(db *sql.DB, testName string) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path.Join(config.CommonConfig.AppDir, "internal/mysqlrepo/fixtures", testName)), // file://path/to/directory
		"mysql", driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")

	db.Close()
}
