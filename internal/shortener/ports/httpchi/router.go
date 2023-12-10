package httpchi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func New(h Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)

	r.Get("/{alias}", h.GetOrigURL)
	r.Post("/alias", h.CrateAlias)
	return r
}
