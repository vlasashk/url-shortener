package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type ShortenerCfg struct {
	App      AppCfg
	Postgres PostgresCfg
}

type AppCfg struct {
	Host      string `env:"APP_HOST" env-default:"localhost"`
	Port      string `env:"APP_PORT" env-default:"9090"`
	BaseURL   string `env:"BASE_URL" env-default:"localhost:9090"`
	LoggerLVL string `env:"APP_LOG_LVL" env-default:"debug"`
}

type PostgresCfg struct {
	Username   string `env:"POSTGRES_USER" env-default:"postgres"`
	Password   string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Port       string `env:"PG_PORT" env-default:"5432"`
	Host       string `env:"POSTGRES_HOST" env-default:"localhost"`
	NameDB     string `env:"POSTGRES_DB" env-default:"postgres"`
	Migrations string `env:"DB_MIGRATION_PATH" env-default:"./migrations"`
}

func NewShortener() (ShortenerCfg, error) {
	var newConfig ShortenerCfg
	if err := cleanenv.ReadEnv(&newConfig); err != nil {
		return ShortenerCfg{}, err
	}
	return newConfig, nil
}
