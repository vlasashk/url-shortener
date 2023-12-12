package httpchi

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResp struct {
	Error string `json:"error"`
}

type AliasResp struct {
	Alias string `json:"alias"`
}

type URLRequest struct {
	Original string `json:"original" validate:"required,url"`
}

func NewErr(err string) ErrResp {
	return ErrResp{
		Error: err,
	}
}

func (resp ErrResp) Send(w http.ResponseWriter, r *http.Request, status int) {
	render.Status(r, status)
	render.JSON(w, r, resp)
}
