package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *API) RegisterRoutes() {
	// register middlewares
	a.mux.Use(middleware.Recoverer)
	a.mux.Use(a.WithLogger())
	a.mux.Use(a.WithSourceService())

	a.mux.Route("/", func(r chi.Router) {
		r.Get("/ping", WithResponse(a.Ping))

		r.Route("/sources", func(r chi.Router) {
			r.Post("/", WithResponse(a.NewSource))
			r.Post("/fetch/{kind}/{id}", WithResponse(a.Fetch))
		})
	})
}
