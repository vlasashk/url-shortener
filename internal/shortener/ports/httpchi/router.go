package httpchi

import (
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vlasashk/url-shortener/config"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/handlers/createurl"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/handlers/origurl"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/middlewares"
	"github.com/vlasashk/url-shortener/internal/shortener/resources"
)

func New(cfg config.Config, res resources.Resources) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middlewares.LoggingMiddleware(res.Log))
	r.Use(middleware.URLFormat)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)

	r.Get("/{alias}", origurl.New(res.Log, res.UseCase).ServeHTTP)
	r.Post("/alias", createurl.New(res.Log, res.UseCase).ServeHTTP)

	return &http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: r,
	}
}
