package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	StoragePath string `env:"DB_URL" env-required:"true"`
	Address     string `env:"SERVER_ADDRESS" env-default:":8080"`
}

func InitConfig() (*Config, error) {
	path := ".env"
	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
