package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/endpoint"
	"github.com/vespaiach/auth/internal/model"
	"github.com/vespaiach/auth/internal/mysqlrepo"
	"github.com/vespaiach/auth/internal/service"
)

func main() {
	var (
		appRepo   *model.AppRepo
		appConfig *conf.AppConfig
		us        service.UserService
		db        *gorm.DB
	)

	appConfig = conf.LoadAppConfig()

	db, err := gorm.Open("mysql", appConfig.DbConfig.BuildMysqlDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	appRepo = mysqlrepo.NewMysqlAppRepo(db)
	us = service.NewUserService(appRepo, appConfig)

	mux := http.NewServeMux()
	mux.Handle("/user/v1/", endpoint.MakeHandler(us, appConfig))

	http.Handle("/", mux)
	http.ListenAndServe(":4000", nil)
}
