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
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vespaiach/auth/internal/appconfig"
)

func TestUserRepo(t *testing.T) {

}

func prepareFixtures(appc *appconfig.AppConfig, fixturePath string) {
	dbConfig := appc.DbConfig

	db, err := sql.Open("mysql", *dbconfig.BuildMysqlDSN())
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
		fmt.Sprintf("file://%s", path.Join(appc.CommonConfig.AppDir, fixturePath)), // file://path/to/directory
		"mysql", driver)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	log.Println("Database migrated")
	// actual logic to start your application
	os.Exit(0)
}
