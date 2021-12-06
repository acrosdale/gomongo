package app

import (
	"github.com/labstack/gommon/log"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/controllers"
	"github.com/acrosdale/gomongo/internal/db"
	"github.com/acrosdale/gomongo/internal/services"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo     *echo.Echo
	Settings *configs.Settings
	Dbs      *db.DBHandler
}

// GetAppServer return a server struct
func GetAppServer() *Server {

	cfg := configs.GetSettings()

	// connect dbs handler
	dbs, err := db.NewDBs(cfg)

	if err != nil {
		log.Fatalf("Unable to conn to db")
		panic(err)
	}

	// create echo routing
	_echo := configs.GetConfiguredEcho()

	// create service layer
	serviceHandler, err := services.CreateServiceHandler(dbs)
	if err != nil {
		log.Fatalf("services unresponsive")
		panic(err)
	}

	// create controllers.. pass service layer into controller
	_, err = controllers.NewController(_echo, serviceHandler)
	if err != nil {
		log.Fatalf("controller unresponsive")
		panic(err)
	}

	return &Server{
		Echo:     _echo,
		Dbs:      dbs,
		Settings: &cfg,
	}
}
