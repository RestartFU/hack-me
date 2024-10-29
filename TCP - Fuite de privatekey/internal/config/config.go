package config

type Config struct {
	Network NetworkConfig `toml:"network"`
}

type NetworkConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func DefaultConfig() Config {
	return Config{
		Network: NetworkConfig{
			Host: "127.0.0.1",
			Port: 8762,
		},
	}
}
