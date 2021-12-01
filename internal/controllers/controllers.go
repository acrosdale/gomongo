/*
	DESCRIPTION
		A controller act as the handler for incoming request
	ERROR RULES:
		A handler MUST ALWAYS return a err DEFINED by the app/service.
		NEVER return error thrown by any go.mod dependencies. for example
		never return a error thrown by mongo to the user (security issue)
*/
package controllers

import (
	apicontroller "github.com/acrosdale/gomongo/internal/controllers/api"
	"github.com/acrosdale/gomongo/internal/services"
	"github.com/labstack/echo/v4"
)

type (
	ControllerHandler struct {
		ApiControllerHandler apicontroller.ApiHandler
	}
)

/*
	NewController return a handler that abstracts the controller layer
*/
func NewController(router *echo.Echo, service *services.ServiceHandler) (*ControllerHandler, error) {
	var err error
	handler := ControllerHandler{}

	handler.ApiControllerHandler, err = apicontroller.NewApiHandler(router, service)

	if err != nil {
		return &ControllerHandler{}, err
	}

	return &handler, nil
}
