package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/migrate"
	"github.com/vespaiach/auth/internal/mysqlrepo"
	"github.com/vespaiach/auth/internal/service"
	tp "github.com/vespaiach/auth/internal/transport"
)

func main() {
	var runmig = flag.Bool("mig", false, "automatically run migration's scripts")
	var adminUsername = flag.String("name", "admin", "administrator's username")
	flag.Parse()

	appConfig := conf.LoadAppConfig()

	db, err := initDb(appConfig)
	if err != nil {
		log.Error("Open database error: ", err)
	}
	defer db.Close()

	logger := log.New()

	if *runmig {
		mig := migrate.NewMigrator(db)
		mig.Down()
		mig.Up()
		mig.SeedProdData("Administrator", *adminUsername, "admin@admin.com", "password", appConfig.BcryptCost)
	}

	appModel := mysqlrepo.NewMysqlAppRepo(db)
	appService := service.NewAppService(appModel, appConfig)

	mux := mux.NewRouter()
	tp.MakeUsersHandlers(mux, appService, appConfig, logger)
	tp.MakeActionHandlers(mux, appService, appConfig, logger)
	tp.MakeRoleHandlers(mux, appService, appConfig, logger)
	tp.MakeUserRoleHandlers(mux, appService, appConfig, logger)
	tp.MakeRoleActionHandlers(mux, appService, appConfig, logger)
	http.Handle("/", mux)

	fmt.Println("transport tp address", appConfig.ServerAddress, "msg listening")
	fmt.Println(http.ListenAndServe(appConfig.ServerAddress, nil))
}

func initDb(config *conf.AppConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", config.BuildMysqlDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
