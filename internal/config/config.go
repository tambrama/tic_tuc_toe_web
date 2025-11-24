package config

type Config struct {
	ServerPort string `env:"SERVER_PORT" default:"8081"`
}

func NewConfig() *Config {
	return &Config{
		ServerPort: "8081",
	}
}
