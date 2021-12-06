package db

import (
	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/acrosdale/gomongo/internal/db/pgdb"
	"github.com/labstack/gommon/log"
)

type (

	// db handler
	DBHandler struct {
		MongoHandler    mgdb.MongoHandler
		PostgresHandler pgdb.PostgresHandler
	}
)

// CloseAllDBConn closes the DBHandler db connection
func CloseAllDBConn(dbs *DBHandler) {
	// disconnect all db clients
	dbs.MongoHandler.CloseMongoHandler()
	dbs.PostgresHandler.ClosePostgresHandler()
}

// NewDBs initd all the dbs can be used by this app
func NewDBs(cfg configs.Settings) (*DBHandler, error) {

	// create mongo handler
	mongoHandler, err := mgdb.NewMongoHandler(cfg.DbConfig.Mgdb)
	if err != nil {
		log.Fatal(err)
	}

	// create postgres handler
	postgresHandler, err := pgdb.NewPostgresHandler(cfg.DbConfig.Pgdb)
	if err != nil {
		log.Fatal(err)
	}

	dbs := &DBHandler{
		MongoHandler:    mongoHandler,
		PostgresHandler: postgresHandler,
	}

	// TODO: ping ALL DBS for verification of connectibilty

	return dbs, nil
}
