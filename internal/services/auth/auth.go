package authservice

import (
	"github.com/acrosdale/gomongo/internal/db"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"golang.org/x/crypto/bcrypt"
)

type (
	// the Product service struct and its requried state
	AuthService struct {
		db *db.DBHandler
	}

	AuthServiceInterface interface {
		UserAuthServicesInterface
	}
)

// CreateAuthService boostrap the services-AUTH layer
func CreateAuthService(db *db.DBHandler) (AuthServiceInterface, error) {
	return AuthService{
		db: db,
	}, nil
}

/* isCredValid compares stored user encryted password(from db) against
the given user passward(from request)
*/
func isCredValid(givenuser mgdb.User, storedUser mgdb.User) bool {

	// create hash from the request user cred and cmp against stored user
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(givenuser.Password))
	return err == nil

}
