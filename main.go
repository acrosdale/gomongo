package main

import (
	"fmt"

	"github.com/acrosdale/gomongo/internal/cmd"
	"github.com/acrosdale/gomongo/internal/db"
)

func main() {

	// get app server
	Server := cmd.GetAppServer()

	// close server...defer
	defer db.CloseAllDBConn(Server.Dbs)

	loggerInitialMsg := fmt.Sprintf(
		"Listening on %s:%s",
		Server.Settings.AppConfig.Host,
		Server.Settings.AppConfig.Port,
	)
	// start app
	Server.Echo.Logger.Infof(loggerInitialMsg)

	Server.Echo.Logger.Infof(fmt.Sprintf("settings: %v", Server.Settings))

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
