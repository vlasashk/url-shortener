package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppCfg
	Postgres PostgresCfg
}

type AppCfg struct {
	Host      string `env:"APP_HOST" env-default:"localhost"`
	Port      string `env:"APP_PORT" env-default:"9090"`
	Address   string `env:"APP_ADDRESS" env-default:"localhost:9090"`
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

func New() (Config, error) {
	var newConfig Config
	if err := cleanenv.ReadEnv(&newConfig); err != nil {
		return Config{}, err
	}
	return newConfig, nil
}
