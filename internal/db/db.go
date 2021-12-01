package db

import (
	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/labstack/gommon/log"
)

type (

	// db handler
	DBHandler struct {
		MongoHandler mgdb.MongoHandler
	}
)

// CloseAllDBConn closes the DBHandler db connection
func CloseAllDBConn(dbs *DBHandler) {
	// disconnect all db clients
	mgdb.CloseMongoHandler(dbs.MongoHandler)

}

// NewDBs initd all the dbs can be used by this app
func NewDBs(cfg configs.Settings) (*DBHandler, error) {

	// create mongo handler
	mongoHandler, err := mgdb.NewMongoHandler(cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
	}

	dbs := &DBHandler{
		MongoHandler: mongoHandler,
	}

	// TODO: ping ALL DBS for verification of connectibilty

	return dbs, nil
}
