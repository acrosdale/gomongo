package apiv1

import (
	"github.com/acrosdale/gomongo/internal/services"
	"github.com/labstack/echo/v4"
)

const apiVersion = "v1"

type (
	Apiv1Handler struct {
		ServiceHandler *services.ServiceHandler
	}

	// apis v1 interface
	ApiInterfaceV1 interface {
		productsAPIs
		userAPIs
		authAPIs
	}
)

// RegisterRoutesV1 binds the v1 endpoints to the router
func RegisterRoutesV1(apis *echo.Group, service *services.ServiceHandler) ApiInterfaceV1 {
	var apiv1handler ApiInterfaceV1 = Apiv1Handler{
		ServiceHandler: service,
	}

	// api versions
	apisv1 := apis.Group("/" + apiVersion)
	bindProductsRoutingV1(apisv1, apiv1handler)
	bindUserRoutingV1(apisv1, apiv1handler)
	bindAuthRoutingV1(apisv1, apiv1handler)
	return apiv1handler
}
