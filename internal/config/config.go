package config

import "github.com/ilyakaznacheev/cleanenv"

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-default:"postgresql"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Pass     string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Database string `env:"POSTGRES_DB" env-default:"postgres"`
}

type App struct {
	Name     string `env:"APP_NAME" env-default:"go-rest-api"`
	LogLevel string `env:"APP_LOG_LEVEL" env-default:"debug"`
}

type Config struct {
	App      App
	Postgres Postgres
}

// New returns a new Config struct
func New() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
