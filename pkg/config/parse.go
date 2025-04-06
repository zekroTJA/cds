package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/file"
)

func Parse(cfgFile string, envPrefix string, def Config) (Config, error) {
	cfg := def

	err := file.Decode(cfgFile, &cfg)
	if err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	_ = godotenv.Load()
	err = env.Decode(os.Environ(), envPrefix, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
