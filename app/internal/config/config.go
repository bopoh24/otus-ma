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

type Keycloak struct {
	Realm        string `env:"KEYCLOAK_REALM" env-default:"app"`
	URL          string `env:"KEYCLOAK_URL" env-default:"http://auth-server-keycloak.auth.svc.cluster.local"`
	ClientID     string `env:"KEYCLOAK_CLIENT_ID" env-default:"simple-server"`
	ClientSecret string `env:"KEYCLOAK_CLIENT_SECRET" env-default:"e2e3f4d5-6c7b-8a9b-0c1d-2e3f4d5e6f7a"`
	Admin        string `env:"KEYCLOAK_ADMIN" env-default:"admin"`
	Password     string `env:"KEYCLOAK_PASSWORD" env-default:"admin"`
}

type Config struct {
	App      App
	Postgres Postgres
	Keycloak Keycloak
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
