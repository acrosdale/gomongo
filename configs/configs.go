package configs

import (
	"github.com/labstack/gommon/log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
	support log format type
*/
const LOG_FORMAT = `${time_rfc3339_nano} ${remote_ip} ${host} ${method} ${uri} ${user_agent}` +
	`${status} ${error} ${latency_human}` + "\n"

const LOG_FORMAT_XCOR_ID = `${time_rfc3339_nano} ${remote_ip} ${header:X-Correlation-ID} ${host} ${method} ${uri} ${user_agent}` +
	`${status} ${error} ${latency_human}` + "\n"

// settings sub struct for a db
type Mgdb struct {
	DBHost string `env:"DB_HOST" env-default:"localhost"`
	DBPort string `env:"DB_PORT" env-default:"-"`
	DBName string `env:"DB_NAME" env-default:"gomongo"`
	DBUser string `env:"DB_USER" env-default:"root"`
	DBPass string `env:"DB_PASS" env-default:"example"`
}
type Pgdb struct {
	DBHost string `env:"POSTGRES_HOST" env-default:"db"`
	DBPort string `env:"POSTGRES_PORT" env-default:"5432"`
	DBName string `env:"POSTGRES_DB" env-default:"gomongo"`
	DBUser string `env:"POSTGRES_USER" env-default:"root"`
	DBPass string `env:"POSTGRES_PASSWORD" env-default:"example"`
}
type DbConfig struct {
	Mgdb Mgdb
	Pgdb Pgdb
}

// a settings sub struct for app related config
type AppConfig struct {
	Port      string `env:"APP_PORT" env-default:"8080"`
	Host      string `env:"HOST" env-default:"gomongo"`
	APPSecret string `env:"APP_SECRET" env-default:"-"`
	AppMode   string `env:"APP_MODE" env-default:"debug"`
}

// the settings config struct
type Settings struct {
	DbConfig  DbConfig
	AppConfig AppConfig
}

// GetSettings return the app setting object
func GetSettings() Settings {
	var cfg Settings

	// load env vars
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Unable to load configurations...in middlewares")
	}

	return cfg
}

// GetConfiguredEcho returns an echo isntance that is configured with basic middlewares
func GetConfiguredEcho() *echo.Echo {

	// create echo routing
	_echo := echo.New()

	// add loggers and middlewares
	_echo.Logger.SetLevel(log.DEBUG)
	_echo.Use(middleware.Recover())
	_echo.Pre(middleware.RemoveTrailingSlash())
	_echo.Pre(AddCorrelationID)
	_echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: LOG_FORMAT_XCOR_ID}))

	return _echo
}
