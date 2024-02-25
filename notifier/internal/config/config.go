package config

import "github.com/ilyakaznacheev/cleanenv"

type App struct {
	Name     string `env:"APP_NAME" env-default:"booking-service"`
	LogLevel string `env:"APP_LOG_LEVEL" env-default:"debug"`
}

type SMTP struct {
	Host     string `env:"SMTP_HOST" env-default:"booksvc-mailhog.booksvc.svc.cluster.local"`
	Port     int    `env:"SMTP_PORT" env-default:"1025"`
	Username string `env:"SMTP_USERNAME" env-default:""`
	Password string `env:"SMTP_PASSWORD" env-default:""`
	From     string `env:"SMTP_FROM" env-default:"no-reply@booksvc.com"`
}

type Kafka struct {
	Hosts string `env:"KAFKA_HOSTS" env-default:"booksvc-kafka-controller-headless.booksvc:9092,booksvc-kafka-controller-headless.booksvc:9093,booksvc-kafka-controller-headless.booksvc:9094"`
}

type Config struct {
	App   App
	SMTP  SMTP
	Kafka Kafka
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
