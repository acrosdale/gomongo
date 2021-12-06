package mgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/acrosdale/gomongo/configs"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	// global interface to store all collection query
	MongoQueries interface {
		ProductQueries
		UserQueries
	}
	// implements the MongoQueries
	MongoQuery struct {
		Mgdb *mongo.Database
	}

	// db client mongo
	MongoHandler struct {
		mongoClient  *mongo.Client // this is privately used to close conn to mongo
		Queryhandler MongoQueries
	}
)

var (
	ErrMongoDB    = errors.New("resource access failed try again")
	ErrInvalidId  = errors.New("invalid id")
	ErrNoIdPassed = errors.New("no id provided")
)

// NewMongoHandler return a NewMongoHandler given the db configs. NewMongoHandler abstact the mongo db layer
func NewMongoHandler(cfg configs.Mgdb) (MongoHandler, error) {
	// create the mongo client
	mongoclient, err := newMongoClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mgdb := mongoclient.Database(cfg.DBName)

	// call collection creation func HERE
	createProductsCollection(mgdb)
	createUserCollection(mgdb)

	var queryhandler MongoQueries = &MongoQuery{
		Mgdb: mgdb,
	}

	return MongoHandler{
		mongoClient:  mongoclient,
		Queryhandler: queryhandler,
	}, nil
}

// newMongoClient return a mongo client given the db configs
func newMongoClient(cfg configs.Mgdb) (*mongo.Client, error) {
	var err error

	// define conn URI
	connectURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
	)

	// init mongo client
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(connectURI))

	if err != nil {
		log.Fatalf("Unable to conn to DB client")
	}

	if err = mongoClient.Connect(context.Background()); err != nil {
		log.Fatalf("Unable to conn to DB")
	}

	// ping db to ensure conn
	if err = mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("Unable to ping to DB")
	}

	return mongoClient, nil
}

// CloseMongoHandler closes the mongo client define in the MongoHandler
func (mongo_handler *MongoHandler) CloseMongoHandler() (bool, error) {
	err := mongo_handler.mongoClient.Disconnect(context.Background())

	if err != nil {
		return false, err
	}

	return true, nil
}
