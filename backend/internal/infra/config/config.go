package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DevMode         bool          `env:"DEV_MODE,default=0"`
	HTTPPort        int           `env:"HTTP_PORT,default=8080"`
	AuthTokenSecret string        `env:"AUTH_TOKEN_SECRET,default=secret"`
	AuthTokenTTL    time.Duration `env:"AUTH_TOKEN_TTL,default=240m"`
}

func Provide(ctx context.Context) Config {
	var c Config

	if err := envconfig.Process(ctx, &c); err != nil {
		panic(err)
	}

	return c
}
