// +build integration

package apiservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
********************* Setup *************************
 */

type TestProductsApiServiceSuiteIntegration struct {
	suite.Suite
	cfg          configs.Settings
	ctx_bkg      context.Context
	testproduct  mgdb.Product
	apiservices  ApiServiceInterface
	filter       map[string]interface{}
	dbs          *db.DBHandler
	singleseeder func(db mgdb.MongoHandler, product mgdb.Product)
}

func (suite *TestProductsApiServiceSuiteIntegration) SetupSuite() {
	suite.cfg = configs.GetSettings()
	suite.ctx_bkg = context.Background()
	suite.singleseeder = mgdb.SeedProductCollectionOne

	// test product
	id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	suite.filter = map[string]interface{}{"_id": id}
	suite.testproduct = mgdb.GetOneProductObj(id)

	var err error
	suite.dbs, err = db.NewDBs(suite.cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("TestProductsApiSuiteIntegration::NewDBs suite setup failed %s", err))
	}

	suite.apiservices, err = CreateApiService(suite.dbs)
	if err != nil {
		log.Fatal(fmt.Sprintf("TestProductsApiSuiteIntegration::CreateApiService suite setup failed %s", err))
	}

	// seed a single entry for look up
	suite.singleseeder(suite.dbs.MongoHandler, suite.testproduct)
}

func (suite *TestProductsApiServiceSuiteIntegration) TearDownSuite() {
	// delete added product
	suite.dbs.MongoHandler.Queryhandler.DeleteOneProduct(suite.ctx_bkg, suite.filter)

	//close db conn
	db.CloseAllDBConn(suite.dbs)
}

// run test suite
func TestProductsApiSuite_Integration(t *testing.T) {
	suite.Run(t, new(TestProductsApiServiceSuiteIntegration))
}

/*
************************* TEST CASES ****************************
 */

func (suite *TestProductsApiServiceSuiteIntegration) TestGetProduct_NoErr() {
	// call service
	result, err := suite.apiservices.GetProduct(suite.ctx_bkg, suite.filter)

	// check return
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.testproduct, result)
}

func (suite *TestProductsApiServiceSuiteIntegration) TestGetProduct_NoFilterPassed() {
	// call service
	result, err := suite.apiservices.GetProduct(suite.ctx_bkg, make(map[string]interface{}))

	// check return
	assert.Equal(suite.T(), mgdb.ErrNoProductFound, err)
	assert.Equal(suite.T(), mgdb.Product{}, result)
}

func (suite *TestProductsApiServiceSuiteIntegration) TestGetProduct_Err() {
	// call service
	result, err := suite.apiservices.GetProduct(suite.ctx_bkg, map[string]interface{}{"_id": primitive.NewObjectID()})

	// check return
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), mgdb.Product{}, result)
}
