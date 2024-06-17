package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type CronCfg struct {
	Schedule  string `env:"CRON_SCHEDULE" env-default:"0 0 * * *"`
	LoggerLVL string `env:"CRON_LOG_LVL" env-default:"debug"`
	Postgres  PostgresCfg
}

func NewCron() (CronCfg, error) {
	var newConfig CronCfg
	if err := cleanenv.ReadEnv(&newConfig); err != nil {
		return CronCfg{}, err
	}
	return newConfig, nil
}
