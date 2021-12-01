package api

import (
	apiv1 "github.com/acrosdale/gomongo/internal/controllers/api/v1"
	"github.com/acrosdale/gomongo/internal/services"
	"github.com/labstack/echo/v4"
)

type (
	ApiHandler struct {
		V1Handler apiv1.ApiInterfaceV1
	}
)

// NewApiHandler return a handler that abstracts the api layer with the controller
func NewApiHandler(router *echo.Echo, service *services.ServiceHandler) (ApiHandler, error) {
	apis := router.Group("/api")
	v1 := apiv1.RegisterRoutesV1(apis, service)

	return ApiHandler{
		V1Handler: v1,
	}, nil
}
