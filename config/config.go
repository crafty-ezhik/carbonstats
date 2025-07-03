package config

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
	Host     string
	Username string
	Password string
}

type LoggerConfig struct {
	Debug bool
}

func LoadConfig() *Config {
	return &Config{}
}
