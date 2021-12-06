package main

import (
	"fmt"

	app "github.com/acrosdale/gomongo/internal/cmd/app"
	"github.com/acrosdale/gomongo/internal/db"
)

func main() {

	// get app server
	Server := app.GetAppServer()

	// close server...defer
	defer db.CloseAllDBConn(Server.Dbs)

	// server start up msg
	Server.Echo.Logger.Infof(fmt.Sprintf(
		"Listening on %s:%s",
		Server.Settings.AppConfig.Host,
		Server.Settings.AppConfig.Port,
	))

	// gomongo map to docker internal host, use ":%s" if mapping is undesired
	Server.Echo.Logger.Fatal(
		Server.Echo.Start(
			fmt.Sprintf(
				":%s",
				Server.Settings.AppConfig.Port,
			),
		),
	)
}
