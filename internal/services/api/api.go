package apiservice

import (
	"github.com/acrosdale/gomongo/internal/db"
)

type (
	ApiServiceInterface interface {
		ProductServiceInterface
	}

	ApiService struct {
		db *db.DBHandler
	}
)

// CreateApiService boostrap the services-API layer
func CreateApiService(db *db.DBHandler) (ApiServiceInterface, error) {
	return ApiService{db: db}, nil
}
