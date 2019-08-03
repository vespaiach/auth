package main

import (
	"fmt"

	"github.com/vespaiach/auth/internal/appconfig"
)

func main() {
	app := appconfig.LoadAppConfig()
	fmt.Println(app)
}
