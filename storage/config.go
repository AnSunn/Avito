package storage

type Config struct {
	//Connection string to DB
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
