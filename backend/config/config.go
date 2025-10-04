package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		HTTP    HTTP
		Auth    Auth
		Log     Log
		PG      PG
		Metrics Metrics
		Swagger Swagger
		WebApi  WebApi
	}

	HTTP struct {
		Port           string `env:"HTTP_PORT" envDefault:"8080"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	Auth struct {
		JwksUrl            string `env:"AUTH_JWKS_URL"`
		NoAuthUserOverride string `env:"NOAUTH_USER_OVERRIDE"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX" envDefault:"2"`
		URL     string `env:"PG_URL,required"`
	}

	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	WebApi struct {
		AerodataboxBaseURL      string `env:"AEDBX_URL" envDefault:"https://aerodatabox.p.rapidapi.com"`
		AerodataboxApiKey       string `env:"AEDBX_APIKEY"`
		DbVendoBaseURL          string `env:"DBVENDO_URL"`
		OpenRouteServiceBaseURL string `env:"ORS_URL" envDefault:"https://api.openrouteservice.org"`
		OpenRouteServiceApiKey  string `env:"ORS_APIKEY"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
