package authservice

import (
	"context"
	"errors"
	"time"

	"github.com/acrosdale/gomongo/configs"
	"github.com/acrosdale/gomongo/internal/db/mgdb"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserAuthServicesInterface interface {
		// AuthenicateUser authenicate a user and return a auth user obj
		AuthenicateUser(ctx context.Context, user mgdb.User) (mgdb.User, error)
		//CreateUser create a new user and return user id hex
		CreateUser(ctx context.Context, user mgdb.User) (string, error)
		//CreateToken create a jwt token for a user
		CreateToken(ctx context.Context, user mgdb.User) (string, error)
	}
)

var (
	// return a when a query to user return no results
	ErrNoUserFound      = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrInvalidCred      = errors.New("credential provided are invalid")
	ErrUserEmailEmpty   = errors.New("user email is empty")
)

func (service AuthService) CreateToken(ctx context.Context, user mgdb.User) (string, error) {
	var cfg configs.Settings = configs.GetSettings()

	if user.Email == "" {
		return "", ErrUserEmailEmpty
	}

	claims := jwt.MapClaims{}
	claims["authorize"] = true
	claims["user_id"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString(([]byte(cfg.AppConfig.APPSecret)))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (service AuthService) CreateUser(ctx context.Context, user mgdb.User) (string, error) {
	var oldUser mgdb.User

	oldUser, err := service.db.MongoHandler.Queryhandler.FindOneUser(ctx, bson.M{"email": user.Email})
	if err != mgdb.ErrNoUserFound && err != nil {
		return "", err
	} else if oldUser.Email == user.Email {
		return "", ErrUserAlreadyExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	if err != nil {
		return "", errors.New("unable to hash password")
	}

	// turn password to hash password
	user.Password = string(hashedPassword)

	insertedID, err := service.db.MongoHandler.Queryhandler.InsertOneUser(ctx, user)

	if err != nil {
		log.Errorf("unable to create user, %v", err)
		return "", err
	}

	return insertedID, nil
}

func (service AuthService) AuthenicateUser(ctx context.Context, user mgdb.User) (mgdb.User, error) {
	var storedUser mgdb.User

	storedUser, err := service.db.MongoHandler.Queryhandler.FindOneUser(ctx, bson.M{"email": user.Email})

	if err != nil {
		return storedUser, err
	}

	if !isCredValid(user, storedUser) {
		return storedUser, ErrInvalidCred
	}

	return mgdb.User{Email: storedUser.Email}, nil
}
