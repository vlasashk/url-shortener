package createurl

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/internal/shortener/ports/httpchi/errhandle"
)

//go:generate mockery --name=AliasCreator
type AliasCreator interface {
	CrateAlias(ctx context.Context, url string) (string, error)
}

type Handler struct {
	creator AliasCreator
	log     zerolog.Logger
}

func New(log zerolog.Logger, creator AliasCreator) *Handler {
	return &Handler{
		creator: creator,
		log:     log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newURL urlRequest
	log := *zerolog.Ctx(r.Context())

	if err := render.DecodeJSON(r.Body, &newURL); err != nil {
		log.Error().Err(err).Send()
		errhandle.NewErr("bad JSON").Send(w, r, http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(newURL); err != nil {
		log.Error().Err(err).Send()
		errhandle.NewErr("invalid JSON").Send(w, r, http.StatusUnprocessableEntity)
		return
	}

	alias, err := h.creator.CrateAlias(r.Context(), newURL.Original)
	if err != nil {
		log.Error().Err(err).Send()
		errhandle.NewErr("alias creation fail").Send(w, r, http.StatusInternalServerError)
		return
	}

	log.Debug().Str("alias", alias).Msg("alias created successfully")
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, aliasResp{Alias: alias})
}
