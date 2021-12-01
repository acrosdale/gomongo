package apiservice

import (
	"context"

	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProductServiceInterface interface {
		// GetProduct returns a mgdb.Product based on the filter passed in
		GetProduct(ctx context.Context, filters map[string]interface{}) (mgdb.Product, error)
		// DeleteProduct deletes a mgdb.Product based on the filter passed in, return count
		DeleteProduct(ctx context.Context, filters map[string]interface{}) (int64, error)
		// InsertProduct insert a mgdb.Product return hex id
		InsertProduct(ctx context.Context, document mgdb.Product) (string, error)
		// InsertProduct updates a mgdb.Product return total doc updated
		UpdateProduct(ctx context.Context, filters map[string]interface{}, document mgdb.Product) (int64, error)
	}
)

func (service ApiService) GetProduct(ctx context.Context, filters map[string]interface{}) (mgdb.Product, error) {
	// var product mgdb.Product

	// no filter return no doc found
	if len(filters) == 0 {
		return mgdb.Product{}, mgdb.ErrNoProductFound
	}

	product, err := service.db.MongoHandler.Queryhandler.FindOneProduct(ctx, filters)

	if err != nil {
		return product, err
	}

	return product, nil
}

func (service ApiService) DeleteProduct(ctx context.Context, filters map[string]interface{}) (int64, error) {
	// no filter return no doc found
	if len(filters) == 0 {
		return 0, mgdb.ErrNoProductFound
	}

	deletedCount, err := service.db.MongoHandler.Queryhandler.DeleteOneProduct(ctx, filters)

	if err != nil {
		return deletedCount, err
	}

	return deletedCount, nil
}

func (service ApiService) InsertProduct(ctx context.Context, document mgdb.Product) (string, error) {

	id, err := service.db.MongoHandler.Queryhandler.InsertOneProduct(ctx, document)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (service ApiService) UpdateProduct(ctx context.Context, filters map[string]interface{}, document mgdb.Product) (int64, error) {
	// dont allow id update
	document.ID = primitive.NilObjectID

	updatedCount, err := service.db.MongoHandler.Queryhandler.UpdateOneProduct(ctx, filters, document)

	if err != nil {
		return 0, err
	}

	return updatedCount, nil
}
