package httpchi

import (
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/models/service"
	"net/http"
)

type Handler struct {
	s   service.Service
	log zerolog.Logger
}

func NewHandler(service service.Service, logger zerolog.Logger) Handler {
	return Handler{
		s:   service,
		log: logger,
	}
}

func Run(service service.Service, logger zerolog.Logger, cfg config.AppCfg) error {
	r := New(NewHandler(service, logger))
	logger.Info().Str("address", cfg.Host+":"+cfg.Port).Msg("starting listening")
	return http.ListenAndServe(cfg.Host+":"+cfg.Port, r)
}
