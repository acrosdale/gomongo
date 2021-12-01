package services

import (
	"log"

	db "github.com/acrosdale/gomongo/internal/db"
	apiservice "github.com/acrosdale/gomongo/internal/services/api"
	authservice "github.com/acrosdale/gomongo/internal/services/auth"
)

/*
	this packages will serves as the encapsulation of business logic that can be
	consume.
*/

type (

	// test
	ServiceHandler struct {
		ApiServices  apiservice.ApiServiceInterface
		AuthServices authservice.AuthServiceInterface
	}
)

// CreateServiceHandler boostrap the service layer and return the ServiceHandler
func CreateServiceHandler(db *db.DBHandler) (*ServiceHandler, error) {
	var err error
	handler := &ServiceHandler{}

	// bind all viable serive handler
	handler.ApiServices, err = apiservice.CreateApiService(db)
	if err != nil {
		log.Fatal(err)
	}

	handler.AuthServices, err = authservice.CreateAuthService(db)
	if err != nil {
		log.Fatal(err)
	}

	return handler, nil
}
