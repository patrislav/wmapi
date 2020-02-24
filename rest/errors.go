package rest

import (
	"github.com/go-chi/render"
	"net/http"
)

type errPayload struct {
	statusCode int `json:"-"`

	Message string `json:"message,omitempty"`
}

func (e errPayload) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.statusCode)
	return nil
}
