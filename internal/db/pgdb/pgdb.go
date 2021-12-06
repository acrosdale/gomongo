package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/acrosdale/gomongo/configs"
	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	// global interface to store all collection query
	PostgresQueries interface {
		// data into table
		InsertComments(ctx context.Context, data []Comment) (int64, error)
		InsertPolls(ctx context.Context, data []Poll) (int64, error)
	}

	// implements the MongoQueries
	PostgresQuery struct {
		Pgdb *bun.DB
	}

	// db client mongo
	PostgresHandler struct {
		client       *bun.DB // this is privately used to close conn to mongo
		Queryhandler PostgresQueries
	}
)

func NewPostgresHandler(cfg configs.Pgdb) (PostgresHandler, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort)),
		pgdriver.WithUser(cfg.DBUser),
		pgdriver.WithPassword(cfg.DBPass),
		pgdriver.WithDatabase(cfg.DBName),
		pgdriver.WithApplicationName("myapp"),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),

		//security
		// pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: false}),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to ping to DB")
	}

	var queryhandler PostgresQueries = &PostgresQuery{Pgdb: db}

	return PostgresHandler{
		client:       db,
		Queryhandler: queryhandler,
	}, nil
}

// CloseMongoHandler closes the mongo client define in the MongoHandler
func (handler *PostgresHandler) ClosePostgresHandler() (bool, error) {
	err := handler.client.Close()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (pghandler *PostgresHandler) GetDatabase() *bun.DB {
	return pghandler.client
}
