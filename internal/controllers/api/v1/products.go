package apiv1

import (
	"fmt"
	"net/http"

	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	productsAPIs interface {
		// getProductAPI retrieve a mgdb.Product from the :id param and return it as a json response
		getProductAPI(ctx echo.Context) error
		// insertProductAPI stores a mgdb.Product and returns the ID
		insertProductAPI(ctx echo.Context) error
		// deleteProductAPI take an :id and deletes the mgdb.Product
		deleteProductAPI(ctx echo.Context) error
		// updateProductAPI take an :id with json product and update mgdb.Product
		updateProductAPI(ctx echo.Context) error
	}
)

// bindProductsRoutingV1 binds the v1 products endpoints to the router
func bindProductsRoutingV1(apiV1Group *echo.Group, handler ApiInterfaceV1) {
	// create apiv1 handler and pass in states
	productsAPI := apiV1Group.Group("/products")

	// use the jwt middleware
	// productsAPI.Use(configs.JwtMiddleware)

	// PRODUCTS ENDPOINTS
	productsAPI.GET("/:id", handler.getProductAPI)
	productsAPI.POST("", handler.insertProductAPI)
	productsAPI.DELETE("/:id", handler.deleteProductAPI)
	productsAPI.PUT("/:id", handler.updateProductAPI)
}

func (handler Apiv1Handler) getProductAPI(ctx echo.Context) error {
	filters := make(map[string]interface{})
	response := make(map[string]interface{})

	// validate and convert id
	if ctx.Param("id") != "" {
		idObj, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		// cannot convert id
		if err != nil {
			response["error"] = mgdb.ErrInvalidId.Error()
			return ctx.JSON(http.StatusBadRequest, response)
		}
		filters["_id"] = idObj
	} else {
		//no id provided
		response["error"] = mgdb.ErrNoIdPassed.Error()
		return ctx.JSON(http.StatusBadRequest, response)
	}

	product, err := handler.ServiceHandler.ApiServices.GetProduct(ctx.Request().Context(), filters)

	// write err to response
	if err != nil {
		response["error"] = err.Error()
		if err == mgdb.ErrMongoDB { // mongo err
			return ctx.JSON(http.StatusInternalServerError, response)

		} else if err == mgdb.ErrNoProductFound { // no products found
			return ctx.JSON(http.StatusNotFound, response)
		}
	}
	response["product"] = product
	return ctx.JSON(http.StatusOK, response)
}

func (handler Apiv1Handler) insertProductAPI(ctx echo.Context) error {
	var product mgdb.Product
	response := make(map[string]interface{})

	// parser request
	if err := ctx.Bind(&product); err != nil {
		log.Errorf("unable to parse data, %v", err)
		response["error"] = "unable to parse request"
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// validate data
	if errs := product.Validate(); errs != nil {
		response["errors"] = errs
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// insert product
	id, err := handler.ServiceHandler.ApiServices.InsertProduct(ctx.Request().Context(), product)

	// determine err code for err
	if err != nil {
		response["error"] = err.Error()
		if err == mgdb.ErrMongoDB { // mongo err
			return ctx.JSON(http.StatusInternalServerError, response)

		}
	}

	response["product"] = id
	return ctx.JSON(http.StatusOK, response)
}

func (handler Apiv1Handler) updateProductAPI(ctx echo.Context) error {
	var product mgdb.Product
	response := make(map[string]interface{})
	filters := make(map[string]interface{})

	// validate and convert id
	if ctx.Param("id") != "" {
		idObj, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		// cannot convert id
		if err != nil {
			response["error"] = mgdb.ErrInvalidId.Error()
			return ctx.JSON(http.StatusBadRequest, response)
		}
		filters["_id"] = idObj
	} else {
		//no id provided
		response["error"] = mgdb.ErrNoIdPassed.Error()
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// parser request
	if err := ctx.Bind(&product); err != nil {
		log.Errorf("unable to parse data, %v", err)
		response["error"] = "unable to parse request"
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// validate data
	if errs := product.Validate(); errs != nil {
		response["errors"] = errs
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// retrienve product from db
	updatedTotal, err := handler.ServiceHandler.ApiServices.UpdateProduct(ctx.Request().Context(), filters, product)

	// write err to response
	if err != nil {
		response["error"] = err.Error()
		if err == mgdb.ErrMongoDB { // mongo err
			return ctx.JSON(http.StatusInternalServerError, response)

		} else if err == mgdb.ErrNoProductFound { // no products found
			return ctx.JSON(http.StatusNotFound, response)
		}
	}
	response["product_updated"] = updatedTotal
	return ctx.JSON(http.StatusOK, response)
}

func (handler Apiv1Handler) deleteProductAPI(ctx echo.Context) error {
	filters := make(map[string]interface{})
	response := make(map[string]interface{})

	// validate and convert id
	if ctx.Param("id") != "" {
		idObj, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		// cannot convert id
		if err != nil {
			response["error"] = mgdb.ErrInvalidId.Error()
			return ctx.JSON(http.StatusBadRequest, response)
		}
		filters["_id"] = idObj
	} else {
		//no id provided
		response["error"] = mgdb.ErrNoIdPassed.Error()
		return ctx.JSON(http.StatusBadRequest, response)
	}

	total, err := handler.ServiceHandler.ApiServices.DeleteProduct(ctx.Request().Context(), filters)

	// write err to response
	if err != nil {
		response["error"] = err.Error()
		if err == mgdb.ErrMongoDB { // mongo err
			return ctx.JSON(http.StatusInternalServerError, response)

		} else if err == mgdb.ErrNoProductFound { // no products found
			return ctx.JSON(http.StatusNotFound, response)
		}
	}

	response["message"] = fmt.Sprintf("%v product(s) deleted", total)
	return ctx.JSON(http.StatusOK, response)
}
