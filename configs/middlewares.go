package configs

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
)

const (
	/*
		this is a request id unique to every request.
		used for request tracing
	*/
	CorrelationID = "X-Correlation-ID"
)

var (
	cfg Settings

	/*
		middlewares
	*/
	JwtMiddleware echo.MiddlewareFunc
)

func init() {

	cfg = GetSettings()

	// setup jwt token middleware
	JwtMiddleware = createJwtMiddleware(cfg)
	middleware.ErrJWTMissing.Code = http.StatusUnauthorized
	middleware.ErrJWTMissing.Message = "Unauthorized"
}

// createJwtMiddleware creates the jwt token to be used in the jwt middleware
func createJwtMiddleware(cfg Settings) echo.MiddlewareFunc {
	jwtmiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(cfg.AppConfig.APPSecret),
		TokenLookup: "header:x-auth-token",
	})

	return jwtmiddleware
}

// AddCorrelationID creates or attacted the X-Correlation-ID for inter-app tracking
func AddCorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var newID string
		id := c.Request().Header.Get(CorrelationID)

		if id == "" {
			// generate correlation ID here
			newID = random.String(12)
		} else {
			newID = id
		}

		c.Request().Header.Set(CorrelationID, newID)
		c.Response().Header().Set(CorrelationID, newID)
		return next(c)
	}
}

// AdminOnlyMiddleware check and verify that the token is authorized as admin to do something
func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// dependent of the JwtMiddleware
	return func(ctx echo.Context) error {
		htoken := ctx.Request().Header.Get("x-auth-token") // Bearer token
		token := strings.Split(htoken, " ")[1]
		var claims = jwt.MapClaims{}

		_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.AppConfig.APPSecret), nil
		})

		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, "unauthorized")
		}

		if !claims["authorized"].(bool) {
			return ctx.JSON(http.StatusForbidden, "Forbidden")
		}
		return next(ctx)
	}
}
