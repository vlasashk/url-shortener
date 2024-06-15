package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = shortener.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}
