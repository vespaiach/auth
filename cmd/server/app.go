package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/adding"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/listing"
	"github.com/vespaiach/auth/pkg/storage/mysql"
	"github.com/vespaiach/auth/pkg/tp"
)

func main() {
	appConfig := cf.LoadAppConfig()

	db, err := mysql.InitDb(appConfig.BuildMysqlDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := mysql.NewStorage(db)

	lstServ := listing.NewService(repo)
	addServ := adding.NewService(repo)

	mux := mux.NewRouter()
	tp.MakeUserHandlers(mux, appConfig, lstServ, addServ)
	http.Handle("/", mux)

	fmt.Println("http address ", appConfig.ServerAddress, " msg listening")
	fmt.Println(http.ListenAndServe(appConfig.ServerAddress, nil))
}
