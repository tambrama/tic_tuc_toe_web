package config

type ConfigDB struct {
	URL string
}

func NewConfigDB() *ConfigDB{
	return &ConfigDB{
		URL: "postgres://postgres:password@localhost:5432/tictuc",
	}
}