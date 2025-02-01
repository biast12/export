package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

type (
	ApiConfig struct {
		PublicKeyPath string `env:"PUBLIC_KEY_PATH" envDefault:"./key.pem.pub"`

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
			Secret        string        `env:"SECRET,required"`
			Expiry        time.Duration `env:"EXPIRY" envDefault:"168h"`
			EncryptionKey string        `env:"ENCRYPTION_KEY,required"`
		} `envPrefix:"JWT_"`

		S3            S3Config            `envPrefix:"S3_"`
		ArtifactStore ArtifactStoreConfig `envPrefix:"ARTIFACT_STORE_"`
	}

	WorkerConfig struct {
		KeyPath string `env:"KEY_PATH" envDefault:"./key.pem"`

		Database DatabaseConfig `envPrefix:"DATABASE_"`

		Daemon struct {
			Interval        time.Duration `env:"INTERVAL" envDefault:"5s"`
			DownloadWorkers int           `env:"DOWNLOAD_WORKERS" envDefault:"250"`
		} `envPrefix:"DAEMON_"`

		S3 S3Config `envPrefix:"S3_"`

		TranscriptS3 struct {
			Buckets       []string `env:"BUCKETS,required"`
			EncryptionKey string   `env:"ENCRYPTION_KEY,required"`
		} `envPrefix:"TRANSCRIPT_S3_"`

		ArtifactStore ArtifactStoreConfig `envPrefix:"ARTIFACT_STORE_"`
	}

	DatabaseConfig struct {
		Uri string `env:"URI,required"`
	}

	S3Config struct {
		AccessKey string `env:"ACCESS_KEY,required"`
		SecretKey string `env:"SECRET_KEY,required"`
		Endpoint  string `env:"ENDPOINT,required"`
		Region    string `env:"REGION,required"`
	}

	ArtifactStoreConfig struct {
		Bucket        string `env:"BUCKET,required"`
		EncryptionKey string `env:"ENCRYPTION_KEY,required"`
	}
)

func New[T any]() (cfg T, err error) {
	err = env.Parse(&cfg)
	return
}
