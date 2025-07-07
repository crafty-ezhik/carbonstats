package routes

import (
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r *chi.Mux, h service_description.ServiceDescriptionHandler) {
	r.Route("/carbonstats", func(r chi.Router) {
		r.Get("/", h.GetByCarbonPK())
	})
}
