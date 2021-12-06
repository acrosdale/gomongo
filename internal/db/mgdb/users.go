package mgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/acrosdale/gomongo/internal/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const usercollection = "users"

type (
	UserQueries interface {
		// FindOneUser find on user in db and return it
		FindOneUser(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (User, error)
		// InsertOneUser creates on user in the db
		InsertOneUser(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (string, error)
	}
	User struct {
		ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Email    string             `json:"email" bson:"email"`
		Password string             `json:"password,omitempty" bson:"password"`
	}
)

var (
	ErrNoUserFound = errors.New("no user found")
)

// Validate, validate the user obj callee and serialize err->map
func (user User) Validate() map[string]interface{} {
	err := validation.ValidateStruct(&user,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&user.Email, validation.Required, is.Email),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&user.Password, validation.Required, validation.Length(8, 300)),
	)
	// serialized the err for response
	return utils.ErrToMap(err)
}

// createUsernameIndex create an index on username in the users collection for faster lookup
func createUsernameIndex(collection *mongo.Collection) {
	uniqeUserIndex := true
	indexField := "email"
	usernameIndex := mongo.IndexModel{
		Keys: bson.D{{Key: indexField, Value: 1}},
		Options: &options.IndexOptions{
			Unique: &uniqeUserIndex,
		},
	}

	_, err := collection.Indexes().CreateOne(context.Background(), usernameIndex)

	if err != nil {
		log.Warn("Unique username index on user collection failed", err)
	}
}

// createUserCollection create the user collection even if unused
func createUserCollection(mgdb *mongo.Database) error {
	collection := mgdb.Collection("users")
	createUsernameIndex(collection)
	return nil
}

func (handler *MongoQuery) FindOneUser(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (User, error) {
	var user User

	result := handler.Mgdb.Collection(usercollection).FindOne(ctx, filter, opts...)

	err := result.Decode(&user)

	// no document found
	if err == mongo.ErrNoDocuments {
		return user, ErrNoUserFound
	} else if err != nil {
		log.Error(fmt.Sprintf("mongodb::UserCollection::FindOne err %v", err))
		return user, ErrMongoDB
	}

	return user, nil
}

func (handler *MongoQuery) InsertOneUser(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (string, error) {
	result, err := handler.Mgdb.Collection(usercollection).InsertOne(ctx, document, opts...)

	if err != nil {
		log.Error(fmt.Sprintf("mongodb::UserCollection::InsertOne err %v", err))
		return "", ErrMongoDB
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
