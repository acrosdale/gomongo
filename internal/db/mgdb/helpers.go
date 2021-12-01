package mgdb

import (
	"context"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SeedProductCollectionOne seed one product to mongo db collection
func SeedProductCollectionOne(db MongoHandler, product Product) {
	_, err := db.Queryhandler.InsertOneProduct(
		context.Background(),
		product,
	)
	if err != nil {
		log.Fatal("mgdb::helpers::SeedProductCollectionOne failed")
	}
}

// GetOneProductObj returns one in-mem mgdb.Product of a given id
func GetOneProductObj(id primitive.ObjectID) Product {
	return Product{
		ID:       id,
		Name:     "test1",
		Price:    250,
		Currency: "JMD",
		Discount: 0,
		Vendor:   "BOC LLC",
		Accessories: []string{
			"charger",
			"gift-coupon",
			"subscription",
		},
	}
}

/*
	USER COLLECTION HELPERS
*/

// GetOneUserObj makes and returns a user obj given the id
func GetOneUserObj(id primitive.ObjectID) User {
	return User{
		ID:       id,
		Email:    "test@gmail.com",
		Password: "password!23",
	}
}
