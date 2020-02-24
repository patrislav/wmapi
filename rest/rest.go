package rest

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	chi.Router
	wm wm
}

func NewServer(wm wm) *Server {
	r := chi.NewRouter()
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Recoverer,
	)

	s := &Server{Router: r, wm: wm}
	r.Get("/windows", s.listAllWindows)
	r.Get("/windows/{windowID}", s.getWindow)
	r.Get("/windows/{windowID}/icon", s.getWindowIcon)

	return s
}
