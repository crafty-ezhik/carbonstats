package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DB     DBConfig
	Carbon CarbonConfig
	Log    LoggerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type CarbonConfig struct {
	Host string
	Port int
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

	return &Config{
		DB: DBConfig{},
		Carbon: CarbonConfig{
			Host: os.Getenv("CARBON_HOST"),
			Port: carbonPort,
		},
		Log: LoggerConfig{
			Debug: os.Getenv("CARBON_DEBUG") == "true",
		},
	}
}
