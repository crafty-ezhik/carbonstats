package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DB     DBConfig
	Carbon CarbonConfig
	Log    LoggerConfig
	Server ServerConfig
}

type ServerConfig struct {
	Port            int
	Timeout         time.Duration
	ShutDownTimeout time.Duration
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type CarbonConfig struct {
	Host    string
	Port    int
	Parents []string
}

type LoggerConfig struct {
	Debug bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	carbonPort, err := strconv.Atoi(os.Getenv("CARBON_PORT"))
	if err != nil {
		carbonPort = 8082
	}

	carbonParents := strings.Split(os.Getenv("CARBON_PARENTS"), ",")

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		serverPort = 8089
	}

	serverTimeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		serverTimeout = 30 * time.Second
	}

	shutdownTimeout, err := time.ParseDuration(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		shutdownTimeout = 5 * time.Second
	}

	return &Config{
		DB: DBConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Username: os.Getenv("DATABASE_USERNAME"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Database: os.Getenv("DATABASE_DATABASE"),
		},
		Carbon: CarbonConfig{
			Host:    os.Getenv("CARBON_HOST"),
			Port:    carbonPort,
			Parents: carbonParents,
		},
		Log: LoggerConfig{
			Debug: os.Getenv("CARBON_DEBUG") == "true",
		},
		Server: ServerConfig{
			Port:            serverPort,
			Timeout:         serverTimeout,
			ShutDownTimeout: shutdownTimeout,
		},
	}
}
