package apiserver

type Config struct {
	BindAddress string `toml:"bind_addr"`
}

func NewConfig() *Config {
	return &Config{
		BindAddress: ":8080",
	}
}
