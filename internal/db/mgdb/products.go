package mgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/acrosdale/gomongo/internal/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const productcollection = "products"

type (
	ProductQueries interface {
		// FindOneProduct find one mgdb.Products and return it
		FindOneProduct(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (Product, error)
		// InsertOneProduct creates one mgdb.Products and return it's ID
		InsertOneProduct(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (string, error)
		// DeleteOneProduct deletes one mgdb.Products and returns DeletedCount
		DeleteOneProduct(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
		// UpdateOneProduct updates one mgdb.Products and returns MatchedCount(updated doc)
		UpdateOneProduct(ctx context.Context, filters interface{}, document interface{}, opts ...*options.UpdateOptions) (int64, error)
	}

	// product schema
	Product struct {
		ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Name        string             `json:"product_name" bson:"product_name" validate:"required,max=10"`
		Price       int                `json:"price" bson:"price" validate:"required,max=2000"`
		Currency    string             `json:"currency" bson:"currency" validate:"required,len=3"`
		Discount    int                `json:"discount" bson:"discount"`
		Vendor      string             `json:"vendor" bson:"vendor" validate:"required"`
		Accessories []string           `json:"accessories,omitempty" bson:"accessories,omitemtpy"`
		IsEssential bool               `json:"is_essential" bson:"is_essential"`
	}
)

var (
	ErrNoProductFound = errors.New("no product found")
)

// Validate, validates the product on that calls it and serialized the err->map
func (product Product) Validate() map[string]interface{} {
	err := validation.ValidateStruct(&product,
		validation.Field(&product.Name, validation.Required, validation.Length(1, 10)),
		validation.Field(&product.Price, validation.Required, validation.Max(2000)),
	)
	// serialized the err for response
	return utils.ErrToMap(err)
}

// createProductsCollection create the products collection even if unused
func createProductsCollection(mgdb *mongo.Database) error {
	mgdb.Collection(productcollection)
	return nil
}

func (handler *MongoQuery) FindOneProduct(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (Product, error) {
	var product Product

	result := handler.Mgdb.Collection(productcollection).FindOne(ctx, filter, opts...)

	err := result.Decode(&product)

	// no document found
	if err == mongo.ErrNoDocuments {
		return product, ErrNoProductFound
	} else if err != nil {
		log.Error(fmt.Sprintf("mongodb::ProductsCollection::FindOne err %v", err))
		return product, ErrMongoDB
	}

	return product, nil
}

func (handler *MongoQuery) InsertOneProduct(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (string, error) {

	result, err := handler.Mgdb.Collection(productcollection).InsertOne(ctx, document, opts...)

	if err != nil {
		log.Error(fmt.Sprintf("mongodb::ProductsCollection::InsertOne err %v", err))
		return "", ErrMongoDB
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (handler *MongoQuery) DeleteOneProduct(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {

	result, err := handler.Mgdb.Collection(productcollection).DeleteOne(ctx, filter, opts...)

	if err != nil {
		log.Error(fmt.Sprintf("mongodb::ProductsCollection::DeleteOne err %v", err))
		return 0, ErrMongoDB
	}

	return result.DeletedCount, nil
}

func (handler *MongoQuery) UpdateOneProduct(ctx context.Context, filters interface{}, document interface{}, opts ...*options.UpdateOptions) (int64, error) {

	update := bson.D{{Key: "$set", Value: document}}
	result, err := handler.Mgdb.Collection(productcollection).UpdateOne(ctx, filters, update, opts...)

	if err != nil {
		log.Error(fmt.Sprintf("mongodb::ProductsCollection::UpdateOne err %v", err))
		return 0, ErrMongoDB
	}

	return result.MatchedCount, nil
}
