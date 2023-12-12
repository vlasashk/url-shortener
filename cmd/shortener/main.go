package main

import (
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/adapters/memrepo"
	"github.com/vlasashk/url-shortener/internal/shortener/adapters/pgrepo"
	"github.com/vlasashk/url-shortener/internal/shortener/models/logger"
	"github.com/vlasashk/url-shortener/internal/shortener/models/service"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi"
)

func main() {
	log := logger.New(zerolog.InfoLevel)
	cfg, err := config.ParseConfigValues()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("config parse success")
	var repo service.Repo
	switch cfg.App.Storage {
	case "postgres":
		pgRepo, err := pgrepo.New(cfg.Postgres)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		defer pgRepo.DB.Close()
		repo = pgRepo
		log.Info().Str("storage", "postgres").Msg("repository init success")
	default:
		repo = memrepo.New()
		log.Info().Str("storage", "memory").Msg("repository init success")
	}

	URLService := service.New(repo, cfg.App.Address)

	log.Fatal().Err(httpchi.Run(URLService, log, cfg.App)).Send()
}
