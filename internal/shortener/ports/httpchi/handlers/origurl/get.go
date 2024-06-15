package origurl

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/errhandle"
)

//go:generate mockery --name=OriginalProvider
type OriginalProvider interface {
	GetOrigURL(ctx context.Context, alias string) (string, error)
}

type Handler struct {
	provider OriginalProvider
	log      zerolog.Logger
}

func New(log zerolog.Logger, provider OriginalProvider) *Handler {
	return &Handler{
		provider: provider,
		log:      log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := *zerolog.Ctx(r.Context())

	alias := chi.URLParam(r, "alias")
	log.Debug().Str("alias", alias).Msg("alias to get")

	originalURL, err := h.provider.GetOrigURL(r.Context(), alias)
	if err != nil {
		log.Error().Err(err).Send()
		errhandle.NewErr("alias search fail").Send(w, r, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
