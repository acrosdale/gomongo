// +build integration

package apiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/acrosdale/gomongo/internal/services"
	"github.com/acrosdale/gomongo/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
********************* Setup *************************
 */

type TestProductsApiSuiteIntegration struct {
	suite.Suite
	cfg              configs.Settings
	ctx_bkg          context.Context
	testproduct      mgdb.Product
	filter           map[string]interface{}
	dbs              *db.DBHandler
	apihandler       ApiInterfaceV1
	request_response map[string]interface{}
	singleseeder     func(db mgdb.MongoHandler, product mgdb.Product)
	echostateGetter  func(method_type string, request_url string) (*http.Request, *httptest.ResponseRecorder, echo.Context)
}

func (suite *TestProductsApiSuiteIntegration) SetupSuite() {
	suite.cfg = configs.GetSettings()
	suite.ctx_bkg = context.Background()
	suite.singleseeder = mgdb.SeedProductCollectionOne
	suite.echostateGetter = utils.CreateEchoTestVars
	suite.request_response = make(map[string]interface{})

	// test product
	id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	suite.filter = map[string]interface{}{"_id": id}
	suite.testproduct = mgdb.GetOneProductObj(id)

	// db setup
	var err error
	suite.dbs, err = db.NewDBs(suite.cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("TestProductsApiSuiteIntegration::NewDBs suite setup failed %s", err))
	}

	// service CreateEchoTestVars
	servicehandler, err := services.CreateServiceHandler(suite.dbs)
	if err != nil {
		log.Fatal(fmt.Sprintf("TestProductsApiSuiteIntegration::CreateServiceHandler suite setup failed %s", err))
	}

	suite.apihandler = RegisterRoutesV1(
		echo.New().Group("/api"),
		servicehandler,
	)

	// seed a single entry for look up
	suite.singleseeder(suite.dbs.MongoHandler, suite.testproduct)
}

func (suite *TestProductsApiSuiteIntegration) TearDownSuite() {
	// delete added product
	suite.dbs.MongoHandler.Queryhandler.DeleteOneProduct(suite.ctx_bkg, suite.filter)

	//close db conn
	db.CloseAllDBConn(suite.dbs)
}

// run test suite
func TestProductsApiSuite_Integration(t *testing.T) {
	suite.Run(t, new(TestProductsApiSuiteIntegration))
}

/*
************************* TEST CASES ****************************
 */

func (suite *TestProductsApiSuiteIntegration) TestGetProduct_NoErr() {
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
	}
}

func (suite *TestProductsApiSuiteIntegration) TestGetProduct_NoProductFound() {
	// get echo states/var
	_, recorder, echo_ctx := suite.echostateGetter(http.MethodGet, "/api/v1/products/:id")

	//add params
	echo_ctx.SetParamNames("id")
	echo_ctx.SetParamValues(primitive.NewObjectID().Hex())

	// run func and check
	if assert.NoError(suite.T(), suite.apihandler.getProductAPI(echo_ctx)) {
		assert.NoError(suite.T(), json.Unmarshal(recorder.Body.Bytes(), &suite.request_response)) // parse data
		assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)                               // check status
		assert.Equal(suite.T(), mgdb.ErrNoProductFound.Error(), suite.request_response["error"])  // check return data
	}
}
