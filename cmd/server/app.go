package main

import (
	"fmt"
	"github.com/vespaiach/auth/pkg/bunchmgr"
	"github.com/vespaiach/auth/pkg/keymgr"
	"github.com/vespaiach/auth/pkg/usrmgr"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vespaiach/auth/pkg/cf"
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

	keyserv := keymgr.NewService(mysql.NewKeyStorage(db))
	bunchserv := bunchmgr.NewService(mysql.NewBunchStorage(db))
	userserv := usrmgr.NewService(mysql.NewUserStorage(db))

	mux := mux.NewRouter()
	tp.MakeKeyHandlers(mux, appConfig, keyserv)
	tp.MakeBunchHandlers(mux, appConfig, bunchserv)
	tp.MakeUserHandlers(mux, appConfig, userserv)
	http.Handle("/", mux)

	fmt.Println("http address ", appConfig.ServerAddress, " msg listening")
	fmt.Println(http.ListenAndServe(appConfig.ServerAddress, nil))
}
