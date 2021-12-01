// +build unit

package apiservice

import (
	"context"
	"testing"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/acrosdale/gomongo/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
********************* Setup *************************
 */

type TestProductApiServiceSuite struct {
	suite.Suite
	cfg         configs.Settings
	ctx_bkg     context.Context
	testproduct mgdb.Product
	MongoMock   *mocks.MongoQueriesMock
	apiservices ApiServiceInterface
	filter      map[string]interface{}
}

func (suite *TestProductApiServiceSuite) SetupSuite() {
	suite.filter = map[string]interface{}{"_id": "507f1f77bcf86cd799439011"}
	suite.cfg = configs.GetSettings()
	suite.ctx_bkg = context.Background()
	// test product
	suite.testproduct = mgdb.GetOneProductObj(primitive.NewObjectID())
}

func (suite *TestProductApiServiceSuite) SetupTest() {
	// mock must be set for each test...else some wierd issue occurs with setting expectation
	suite.MongoMock = new(mocks.MongoQueriesMock)
	suite.apiservices = ApiService{
		db: &db.DBHandler{
			MongoHandler: mgdb.MongoHandler{
				Queryhandler: suite.MongoMock,
			},
		},
	}
}

// run test suite
func TestProductApiService(t *testing.T) {
	suite.Run(t, new(TestProductApiServiceSuite))
}

/*
************************* TEST CASES ****************************
 */

func (suite *TestProductApiServiceSuite) TestGetProduct_NoErr() {
	// set expection on Mock
	suite.MongoMock.
		On("FindOneProduct", suite.ctx_bkg, suite.filter).
		Return(suite.testproduct, nil)

	// call service
	result, err := suite.apiservices.GetProduct(suite.ctx_bkg, suite.filter)

	// check return
	suite.MongoMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.testproduct, result)
}

func (suite *TestProductApiServiceSuite) TestGetProduct_NoProductFound() {
	emptyProducts := mgdb.Product{}
	// set expection on Mock
	suite.MongoMock.
		On("FindOneProduct", suite.ctx_bkg, suite.filter).
		Return(emptyProducts, mgdb.ErrNoProductFound)

	// call service
	_, err := suite.apiservices.GetProduct(suite.ctx_bkg, suite.filter)

	// check return
	suite.MongoMock.AssertExpectations(suite.T())
	assert.Equal(suite.T(), mgdb.ErrNoProductFound, err)
}

func (suite *TestProductApiServiceSuite) TestGetProduct_NoFilterPassed() {
	filter := make(map[string]interface{})

	// call service
	_, err := suite.apiservices.GetProduct(suite.ctx_bkg, filter)

	// check return
	// suite.productcollectionMock.AssertExpectations(suite.T())
	assert.Equal(suite.T(), mgdb.ErrNoProductFound, err)
}

func (suite *TestProductApiServiceSuite) TestDeleteProduct_NoErr() {
	var total_products_deleted int64 = 1
	// set expection on Mock
	suite.MongoMock.
		On("DeleteOneProduct", suite.ctx_bkg, suite.filter).
		Return(total_products_deleted, nil)

	// call service
	result, err := suite.apiservices.DeleteProduct(suite.ctx_bkg, suite.filter)

	// check return
	suite.MongoMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), total_products_deleted, result)
}

func (suite *TestProductApiServiceSuite) TestDeleteProduct_NoFilterPassed() {
	var total_products_deleted int64 = 0
	emty_filter := make(map[string]interface{})
	// call service
	result, err := suite.apiservices.DeleteProduct(suite.ctx_bkg, emty_filter)

	// check return
	assert.Equal(suite.T(), mgdb.ErrNoProductFound, err)
	assert.Equal(suite.T(), total_products_deleted, result)
}

func (suite *TestProductApiServiceSuite) TestDeleteProduct_ErrMongoDB() {
	var total_products_deleted int64 = 0
	// set expection on Mock
	suite.MongoMock.
		On("DeleteOneProduct", suite.ctx_bkg, suite.filter).
		Return(total_products_deleted, mgdb.ErrMongoDB)

	// call service
	result, err := suite.apiservices.DeleteProduct(suite.ctx_bkg, suite.filter)

	// check return
	suite.MongoMock.AssertExpectations(suite.T())
	assert.Equal(suite.T(), mgdb.ErrMongoDB, err)
	assert.Equal(suite.T(), total_products_deleted, result)
}
