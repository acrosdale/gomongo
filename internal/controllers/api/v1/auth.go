package apiv1

import (
	"net/http"

	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/labstack/echo/v4"
)

type (
	authAPIs interface {
		// AuthenicateUserAPI authenicates a user and set the x-auth-token on the response
		AuthenicateUserAPI(ctx echo.Context) error
	}
)

// bindAuthRoutingV1 binds the v1 auth endpoints to the router
func bindAuthRoutingV1(apiV1Group *echo.Group, handler ApiInterfaceV1) {
	authAPIs := apiV1Group.Group("/auth")

	// User ENDPOINTS
	authAPIs.POST("", handler.AuthenicateUserAPI)

}

func (handler Apiv1Handler) AuthenicateUserAPI(ctx echo.Context) error {
	var user mgdb.User
	response := make(map[string]interface{})

	if err := ctx.Bind(&user); err != nil {
		response["error"] = "data provided cannot be processed"
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// validate user
	if errs := user.Validate(); errs != nil {
		response["error"] = errs
		return ctx.JSON(http.StatusBadRequest, response)
	}

	//auth user
	authuser, err := handler.ServiceHandler.AuthServices.AuthenicateUser(ctx.Request().Context(), user)

	if err != nil {
		response["error"] = err.Error()
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	// create token
	token, err := handler.ServiceHandler.AuthServices.CreateToken(ctx.Request().Context(), authuser)

	if err != nil {
		response["error"] = err.Error()
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// // set token on header
	ctx.Response().Header().Set("x-auth-token", "Bearer "+token)
	response["message"] = "user authenicated"
	return ctx.JSON(http.StatusOK, response)
}
