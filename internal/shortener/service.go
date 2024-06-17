package shortener

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi"
	"github.com/vlasashk/url-shortener/internal/shortener/resources"
	"golang.org/x/sync/errgroup"
)

const gracefulTimeout = 10 * time.Second

func Run(ctx context.Context, cfg config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	res, err := resources.NewResources(ctx, cfg)
	if err != nil {
		return err
	}
	defer res.Stop()

	srv := httpchi.New(cfg, res)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info().Msg(fmt.Sprintf("starting server: %s", net.JoinHostPort(cfg.App.Host, cfg.App.Port)))
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Got interruption signal")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
		defer cancel()
		return srv.Shutdown(shutDownCtx)
	})

	if err = g.Wait(); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	log.Info().Msg("server was gracefully shut down")
	return nil
}
