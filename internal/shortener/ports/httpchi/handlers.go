package httpchi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/vlasashk/url-shortener/internal/shortener/models/urlalias"
	"net/http"
)

func (h Handler) CrateAlias(w http.ResponseWriter, r *http.Request) {
	newUrl := urlalias.URL{}
	log := *zerolog.Ctx(r.Context())
	if err := render.DecodeJSON(r.Body, &newUrl); err != nil {
		log.Error().Err(err).Send()
		NewErr("bad JSON").Send(w, r, http.StatusBadRequest)
		return
	}
	log.Info().Msg("request body decoded")
	if err := validator.New().Struct(newUrl); err != nil {
		log.Error().Err(err).Send()
		NewErr("invalid JSON").Send(w, r, http.StatusUnprocessableEntity)
		return
	}
	alias, err := h.s.CrateAlias(newUrl.Original)
	if err != nil {
		log.Error().Err(err).Send()
		NewErr("alias creation fail").Send(w, r, http.StatusInternalServerError)
		return
	}
	log.Info().Msg("alias created successfully")
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, AliasResp{Alias: alias})
}

func (h Handler) GetOrigURL(w http.ResponseWriter, r *http.Request) {
	log := *zerolog.Ctx(r.Context())
	alias := chi.URLParam(r, "alias")
	log.Info().Str("id", alias).Msg("alias received")
	originalURL, err := h.s.GetOrigURL(alias)
	if err != nil {
		log.Error().Err(err).Send()
		NewErr("alias search fail").Send(w, r, http.StatusNotFound)
		return
	}
	log.Info().Str("id", alias).Msg("received successfully")
	http.Redirect(w, r, originalURL, http.StatusFound)
}
