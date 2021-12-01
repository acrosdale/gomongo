package apiv1

import (
	"net/http"

	"github.com/acrosdale/gomongo/internal/db/mgdb"
	authservice "github.com/acrosdale/gomongo/internal/services/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	userAPIs interface {
		// createUserAPI create user
		createUserAPI(ctx echo.Context) error
	}
)

// bindUserRoutingV1 binds the v1 User endpoints to the router
func bindUserRoutingV1(apiV1Group *echo.Group, handler ApiInterfaceV1) {
	usersAPIs := apiV1Group.Group("/users")

	// User ENDPOINTS
	usersAPIs.POST("", handler.createUserAPI)

}

func (handler Apiv1Handler) createUserAPI(ctx echo.Context) error {
	var user mgdb.User
	var response = make(map[string]interface{})

	// parse request
	if err := ctx.Bind(&user); err != nil {
		log.Errorf("parsing err in createUserAPI, %v", err)
		response["error"] = "unable to parse request"
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// validate data
	if errs := user.Validate(); errs != nil {
		response["errors"] = errs
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// create user
	_, err := handler.ServiceHandler.AuthServices.CreateUser(ctx.Request().Context(), user)
	if err != nil {
		// if the user is already created return success
		if authservice.ErrUserAlreadyExist == err {
			response["message"] = "User created successfully"
			return ctx.JSON(http.StatusOK, response)
		}

		log.Errorf("mgdb.User validation failed, %v", err)
		response["error"] = err.Error()
		return ctx.JSON(http.StatusBadRequest, response)
	}
	response["message"] = "User created succesfully"

	return ctx.JSON(http.StatusCreated, response)
}
