package httpchi

import "github.com/vlasashk/url-shortener/internal/shortener/usecase"

type Handler struct {
	service usecase.Service
}

func NewHandler(s usecase.Service) Handler {
	return Handler{
		service: s,
	}
}
