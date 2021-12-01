package apiv1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/acrosdale/gomongo/internal/mocks"
	"github.com/acrosdale/gomongo/internal/services"
	"github.com/acrosdale/gomongo/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestProductAPISuite struct {
	suite.Suite
	cfg              configs.Settings
	ctx_bkg          context.Context
	apihandler       ApiInterfaceV1
	filter           map[string]interface{}
	testproduct      mgdb.Product
	apiServicesMock  *mocks.ApiServicesMock
	echostateGetter  func(method_type string, request_url string) (*http.Request, *httptest.ResponseRecorder, echo.Context)
	request_response map[string]interface{}
}

func (suite *TestProductAPISuite) SetupSuite() {
	suite.filter = map[string]interface{}{"_id": primitive.NewObjectID()}
	suite.cfg = configs.GetSettings()
	suite.ctx_bkg = context.Background()
	suite.echostateGetter = utils.CreateEchoTestVars
	// test product
	suite.testproduct = mgdb.GetOneProductObj(suite.filter["_id"].(primitive.ObjectID))
}

func (suite *TestProductAPISuite) SetupTest() {
	suite.apiServicesMock = new(mocks.ApiServicesMock)
	// create mock services
	suite.apihandler = Apiv1Handler{
		ServiceHandler: &services.ServiceHandler{
			ApiServices: suite.apiServicesMock,
		},
	}
	suite.request_response = make(map[string]interface{})
}

func TestProductAPIS(t *testing.T) {
	suite.Run(t, new(TestProductAPISuite))
}

/*	********************* Test Cases******************	*/
func (suite *TestProductAPISuite) TestProductAPI_NoErr() {
	// set mocked expectation
	suite.apiServicesMock.
		On("GetProduct", suite.ctx_bkg, suite.filter).
		Return(suite.testproduct, nil)

	// get echo states/var
	_, recorder, echo_ctx := suite.echostateGetter(http.MethodGet, "/api/v1/products/:id")

	//add params
	echo_ctx.SetParamNames("id")
	echo_ctx.SetParamValues(suite.testproduct.ID.Hex())

	// run func and check
	if assert.NoError(suite.T(), suite.apihandler.getProductAPI(echo_ctx)) {
		var product mgdb.Product
		assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &suite.request_response)) // parse data
		assert.NoError(suite.T(), utils.ToStruct(suite.request_response["product"], &product), "response converter failed")
		assert.Equal(suite.T(), http.StatusOK, recorder.Code) // check status
		assert.Equal(suite.T(), suite.testproduct, product)   // check return data
		suite.apiServicesMock.AssertExpectations(suite.T())   // assert expectation
	}
}
func (suite *TestProductAPISuite) TestProductAPI_InvalidId() {
	// get echo states/var
	_, recorder, echo_ctx := suite.echostateGetter(http.MethodGet, "/api/v1/products/:id")

	//add params
	echo_ctx.SetParamNames("id")
	echo_ctx.SetParamValues("invalid_id")

	// make api call
	if assert.NoError(suite.T(), suite.apihandler.getProductAPI(echo_ctx)) {
		_ = json.Unmarshal(recorder.Body.Bytes(), &suite.request_response)                  // parse json
		assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)                       // check status
		assert.Equal(suite.T(), mgdb.ErrInvalidId.Error(), suite.request_response["error"]) // check return data
	}
}

func (suite *TestProductAPISuite) TestProductAPI_NoFilter() {
	// get echo states/var
	_, recorder, echo_ctx := suite.echostateGetter(http.MethodGet, "/api/v1/products/:id")

	//add params
	echo_ctx.SetParamNames("id")
	echo_ctx.SetParamValues("")

	// run func and check
	if assert.NoError(suite.T(), suite.apihandler.getProductAPI(echo_ctx)) {
		_ = json.Unmarshal(recorder.Body.Bytes(), &suite.request_response)                   // parse json
		assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)                        // check status
		assert.Equal(suite.T(), mgdb.ErrNoIdPassed.Error(), suite.request_response["error"]) // check return data
	}
}

func (suite *TestProductAPISuite) TestProductAPI_ErrMongoDB() {
	// set mocked expectation
	suite.apiServicesMock.
		On("GetProduct", suite.ctx_bkg, suite.filter).
		Return(mgdb.Product{}, mgdb.ErrMongoDB)

	// get echo states/var
	_, recorder, echo_ctx := suite.echostateGetter(http.MethodGet, "/api/v1/products/:id")

	//add params
	echo_ctx.SetParamNames("id")
	echo_ctx.SetParamValues(suite.filter["_id"].(primitive.ObjectID).Hex())

	// run func and check
	if assert.NoError(suite.T(), suite.apihandler.getProductAPI(echo_ctx)) {
		assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &suite.request_response))

		assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)            // check status
		assert.Equal(suite.T(), mgdb.ErrMongoDB.Error(), suite.request_response["error"]) // checkr return data
		suite.apiServicesMock.AssertExpectations(suite.T())                               // assert expectation
	}
}

func (suite *TestProductAPISuite) TestProductAPI_ErrNoProductFound() {
	// set mocked expectation
	suite.apiServicesMock.
		On("GetProduct", suite.ctx_bkg, suite.filter).
		Return(mgdb.Product{}, mgdb.ErrNoProductFound)
	// get echo states/var
	_, recorder, echo_ctx := suite.echostateGetter(http.MethodGet, "/api/v1/products/:id")

	//add params
	echo_ctx.SetParamNames("id")
	echo_ctx.SetParamValues(suite.filter["_id"].(primitive.ObjectID).Hex())

	// run func and check
	if assert.NoError(suite.T(), suite.apihandler.getProductAPI(echo_ctx)) {
		suite.apiServicesMock.AssertExpectations(suite.T())                                       // assert expectation
		assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &suite.request_response)) // parse response
		assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)                               // check status
		assert.Equal(suite.T(), mgdb.ErrNoProductFound.Error(), suite.request_response["error"])  // checkr return data
	}
}
