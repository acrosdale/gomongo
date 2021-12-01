package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

// CreateEchoTestVars create the necessary artifacts to mock a request
func CreateEchoTestVars(method_type string, request_url string) (*http.Request, *httptest.ResponseRecorder, echo.Context) {
	var request *http.Request
	echo := echo.New()

	if method_type == http.MethodGet {
		request = httptest.NewRequest(http.MethodGet, "/", nil)
	} else {
		panic(fmt.Sprintf("method_type %s not supported", method_type))
	}

	recorder := httptest.NewRecorder()

	echo_ctx := echo.NewContext(request, recorder)

	echo_ctx.SetPath(request_url)

	return request, recorder, echo_ctx

}
