package cronjob

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/cronjob/tasks"
	"github.com/vlasashk/url-shortener/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg config.CronCfg) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, err := logger.New(cfg.LoggerLVL)
	if err != nil {
		return err
	}

	if err = testConnToDB(ctx, cfg.Postgres, log); err != nil {
		return err
	}

	c := cron.New()

	if _, err = c.AddFunc(cfg.Schedule, func() {
		log.Info().Msg("executing cron job")
		conn, err := tasks.New(ctx, cfg.Postgres, log)
		if err != nil {
			log.Error().Err(err).Msg("failed to connect to postgres")
			return
		}
		defer conn.Close()
		if err = conn.CleanUp(ctx); err != nil {
			log.Error().Err(err).Msg("failed to clean up")
		}
		log.Info().Msg("cron job finished")
	}); err != nil {
		return err
	}

	log.Info().Msg("cron job started")
	c.Start()
	defer c.Stop()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Got interruption signal")
		return nil
	})

	if err = g.Wait(); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	log.Info().Msg("cron was gracefully shut down")
	return nil
}

func testConnToDB(ctx context.Context, cfg config.PostgresCfg, log zerolog.Logger) error {
	conn, err := tasks.New(ctx, cfg, log)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
