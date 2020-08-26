package apiserver

import "github.com/Wardenclock1759/StoreAPI/internal/storage"

type Config struct {
	BindAddress string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	Storage     *storage.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddress: ":8080",
		LogLevel:    "debug",
		Storage:     storage.NewConfig(),
	}
}
