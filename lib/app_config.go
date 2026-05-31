package lib

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName       string
	AppEnv        string
	AppDebug      bool
	AppPort       string
	ServerTimeOut int

	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string

	DBSeedDefaultUser bool
	DBRunMigrations   bool

	AuthDefaultUsername string
	AuthDefaultPassword string
	AuthDefaultName     string
}

func NewAppConfig() *AppConfig {
	_ = godotenv.Load()

	serverTimeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))

	appDebug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	dbSeedDefaultUser, _ := strconv.ParseBool(os.Getenv("DB_SEED_DEFAULT_USER"))
	dbRunMigrations, _ := strconv.ParseBool(os.Getenv("DB_RUN_MIGRATIONS"))

	return &AppConfig{
		AppName:       os.Getenv("APP_NAME"),
		AppEnv:        os.Getenv("APP_ENV"),
		AppDebug:      appDebug,
		AppPort:       os.Getenv("APP_PORT"),
		ServerTimeOut: serverTimeout,

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		DBSeedDefaultUser: dbSeedDefaultUser,
		DBRunMigrations:   dbRunMigrations,

		AuthDefaultUsername: os.Getenv("AUTH_DEFAULT_USERNAME"),
		AuthDefaultPassword: os.Getenv("AUTH_DEFAULT_PASSWORD"),
		AuthDefaultName:     os.Getenv("AUTH_DEFAULT_NAME"),
	}
}
