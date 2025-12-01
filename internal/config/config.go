package config

type Config struct {
	ServerPort string `env:"SERVER_PORT" default:"8081"`
}

func NewConfig() *Config {
	return &Config{
		ServerPort: "8081",
	}
}

type ConfigDB struct {
	URL string
}

func NewConfigDB() *ConfigDB {
	return &ConfigDB{
		URL: "postgres://postgres:password@localhost:5432/tictuc",
	}
}
