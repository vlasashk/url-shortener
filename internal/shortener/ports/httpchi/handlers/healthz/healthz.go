package healthz

import (
	"net/http"

	"github.com/go-chi/render"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"status": "ok",
	})
}
