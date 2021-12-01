// +build unit

package authservice

import (
	"context"
	"testing"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/acrosdale/gomongo/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestAuthUserServiceSuite struct {
	suite.Suite
	cfg     configs.Settings
	ctx_bkg context.Context
	// testproduct mgdb.Product
	MongoMock    *mocks.MongoQueriesMock
	authservices AuthService
}

func (suite *TestAuthUserServiceSuite) SetupSuite() {
	// suite.filter = map[string]interface{}{"_id": "507f1f77bcf86cd799439011"}
	suite.cfg = configs.GetSettings()
	suite.ctx_bkg = context.Background()
	// // test product
	// suite.testproduct = mgdb.GetOneProductObj(primitive.NewObjectID())
}

func (suite *TestAuthUserServiceSuite) SetupTest() {
	// mock must be set for each test...else some wierd issue occurs with setting expectation
	suite.MongoMock = new(mocks.MongoQueriesMock)
	suite.authservices = AuthService{
		db: &db.DBHandler{
			MongoHandler: mgdb.MongoHandler{
				Queryhandler: suite.MongoMock,
			},
		},
	}
}

// run test suite
func TestAuthUserService(t *testing.T) {
	suite.Run(t, new(TestAuthUserServiceSuite))
}

/* tests*/

func (suite *TestAuthUserServiceSuite) TestCreateToken_NoErr() {
	user := mgdb.GetOneUserObj(primitive.NewObjectID())

	result_string, err := suite.authservices.CreateToken(suite.ctx_bkg, user)
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), result_string, "")
}

func (suite *TestAuthUserServiceSuite) TestCreateToken_NoEmail() {
	user := mgdb.GetOneUserObj(primitive.NewObjectID())
	user.Email = ""

	result_string, err := suite.authservices.CreateToken(suite.ctx_bkg, user)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), result_string, "")
	assert.Equal(suite.T(), ErrUserEmailEmpty, err)
}

func (suite *TestAuthUserServiceSuite) TestCreateUser_NoErr() {
	user := mgdb.GetOneUserObj(primitive.NewObjectID())

	suite.MongoMock.
		On("FindOneUser", suite.ctx_bkg, bson.M{"email": user.Email}).
		Return(mgdb.User{}, nil)

	suite.MongoMock.
		On("InsertOneUser", suite.ctx_bkg, mock.Anything).
		Return(user.ID.Hex(), nil)

	id, err := suite.authservices.CreateUser(suite.ctx_bkg, user)

	suite.MongoMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), user.ID, id)
}

func (suite *TestAuthUserServiceSuite) TestCreateUser_DuplicateEmail() {
	user := mgdb.GetOneUserObj(primitive.NewObjectID())

	suite.MongoMock.
		On("FindOneUser", suite.ctx_bkg, bson.M{"email": user.Email}).
		Return(mgdb.User{Email: user.Email}, nil)

	_, err := suite.authservices.CreateUser(suite.ctx_bkg, user)

	suite.MongoMock.AssertExpectations(suite.T())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrUserAlreadyExist, err)
}

func (suite *TestAuthUserServiceSuite) TestCreateUser_NoUser() {
	user := mgdb.GetOneUserObj(primitive.NewObjectID())

	suite.MongoMock.
		On("FindOneUser", suite.ctx_bkg, bson.M{"email": user.Email}).
		Return(mgdb.User{}, mgdb.ErrMongoDB)

	_, err := suite.authservices.CreateUser(suite.ctx_bkg, user)

	suite.MongoMock.AssertExpectations(suite.T())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), mgdb.ErrMongoDB, err)
}

func (suite *TestAuthUserServiceSuite) TestCreateUser_ErrInsert() {
	user := mgdb.GetOneUserObj(primitive.NewObjectID())

	suite.MongoMock.
		On("FindOneUser", suite.ctx_bkg, bson.M{"email": user.Email}).
		Return(mgdb.User{}, nil)

	suite.MongoMock.
		On("InsertOneUser", suite.ctx_bkg, mock.Anything).
		Return("", mgdb.ErrMongoDB)

	result, err := suite.authservices.CreateUser(suite.ctx_bkg, user)

	suite.MongoMock.AssertExpectations(suite.T())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), mgdb.ErrMongoDB, err)
	assert.Equal(suite.T(), result, "")
}
