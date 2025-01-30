package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

type (
	ApiConfig struct {
		Server struct {
			Address string `env:"ADDRESS" envDefault:":8080"`
		} `envPrefix:"SERVER_"`

		Database DatabaseConfig `envPrefix:"DATABASE_"`

		Discord struct {
			RootUrl      string `env:"ROOT_URL" envDefault:"https://discord.com"`
			ClientId     string `env:"CLIENT_ID,required"`
			ClientSecret string `env:"CLIENT_SECRET,required"`
			RedirectUri  string `env:"REDIRECT_URI,required"`
		} `envPrefix:"DISCORD_"`

		Jwt struct {
			Secret string        `env:"SECRET,required"`
			Expiry time.Duration `env:"EXPIRY" envDefault:"168h"`
		} `envPrefix:"JWT_"`
	}

	WorkerConfig struct {
		Database DatabaseConfig `envPrefix:"DATABASE_"`
	}

	DatabaseConfig struct {
		Uri string `env:"URI,required"`
	}
)

func New() (cfg ApiConfig, err error) {
	err = env.Parse(&cfg)
	return
}
