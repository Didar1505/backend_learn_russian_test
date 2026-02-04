package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	JWTSecret   string `env:"JWT_SECRET" env-required:"true"`
	// RefreshSecret   string `env:"REFRESH_SECRET" env-required:"true"`
	// RefreshTTLHours int    `env:"REFRESH_TTL_HOURS" env-required:"true"`

	SMTPHost string `env:"SMTP_HOST" env-required:"true"`
	SMTPPort int    `env:"SMTP_PORT" env-required:"true"`
	SMTPUser string `env:"SMTP_USER" env-required:"true"`
	SMTPPass string `env:"SMTP_PASS" env-required:"true"`
	SMTPFrom string `env:"SMTP_FROM" env-required:"true"`

	GoogleOAuthCredentials  string `env:"GOOGLE_OAUTH_CREDENTIALS"`
	GoogleOAuthRedirectURL  string `env:"GOOGLE_OAUTH_REDIRECT_URL"`
	GoogleOAuthCookieSecret string `env:"GOOGLE_OAUTH_COOKIE_SECRET"`
}

func Load() (*Config, error) {
	var cfg Config
	var err error

	if _, statErr := os.Stat(".env"); statErr == nil {
		err = cleanenv.ReadConfig(".env", &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
