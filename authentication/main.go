package main

import (
	"fmt"
	"log"

	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vespaiach/authentication/config"
	"github.com/vespaiach/authentication/logging"
	"github.com/vespaiach/authentication/service"
	"github.com/vespaiach/authentication/transport"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("couldn't load configuration")
	}

	db, err := initDb(config)
	if err != nil {
		log.Fatal("couldn't connect to database")
	}
	defer db.Close()

	logger := createLogger()
	authService := service.NewService(db, config)
	authService = logging.NewLogging(logger, authService)

	mux := http.NewServeMux()
	mux.Handle("/users/v1/", transport.MakeHandler(authService, logger))

	http.ListenAndServe(":8080", mux)
}

func initDb(config *config.ServiceConfig) (*gorm.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", config.DbUser, config.DbPass, config.DbHost,
		config.DbPort, config.DbName, "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0")

	db, err := gorm.Open("mysql", conn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createLogger() kitlog.Logger {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	return logger
}
