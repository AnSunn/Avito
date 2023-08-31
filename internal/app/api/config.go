package api

import "github.com/AnSunn/ServerUserSegmentation/storage"

// General config for rest api
type Config struct {
	//Port for start api
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *storage.Config
}

// Should return default config
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8181",
		LogLevel: "debug",
		Store:    storage.NewConfig(),
	}
}
