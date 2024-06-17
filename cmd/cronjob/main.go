package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/cronjob"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewCron()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = cronjob.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}
