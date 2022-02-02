package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DevMode         bool          `default:"0" split_words:"true"`
	HTTPPort        int           `default:"8080" split_words:"true"`
	AuthTokenSecret string        `default:"secret" split_words:"true"`
	AuthTokenTTL    time.Duration `default:"240m" split_words:"true"`
}

func Provide() Config {
	var c Config

	if err := envconfig.Process("APP", &c); err != nil {
		panic(err)
	}

	return c
}
