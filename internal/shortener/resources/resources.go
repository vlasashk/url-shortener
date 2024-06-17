package resources

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/adapters/pgrepo"
	"github.com/vlasashk/url-shortener/internal/shortener/usecase"
	"github.com/vlasashk/url-shortener/pkg/logger"
)

type Resources struct {
	Log           zerolog.Logger
	UseCase       *usecase.UseCase
	stopResources []func() error
}

func NewResources(ctx context.Context, cfg config.ShortenerCfg) (Resources, error) {
	log, err := logger.New(cfg.App.LoggerLVL)
	if err != nil {
		return Resources{}, err
	}

	pgRepo, err := pgrepo.New(ctx, cfg.Postgres, log)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return Resources{
		Log: log,
		UseCase: usecase.New(
			cfg.App.Address,
			pgRepo,
			pgRepo,
		),
		stopResources: []func() error{
			func() error {
				pgRepo.DB.Close()
				return nil
			},
		},
	}, nil
}

func (r Resources) Stop() {
	for _, stop := range r.stopResources {
		if err := stop(); err != nil {
			r.Log.Error().Err(err).Send()
		}
	}
}
