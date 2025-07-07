package routes

import (
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r *chi.Mux, h service_description.ServiceDescriptionHandler) {
	r.Route("/carbonstats", func(r chi.Router) {
		r.Route("/services_desc", func(r chi.Router) {
			r.Get("/{carbon_pk}", h.GetByCarbonPK())
			r.Get("/", h.GetAll())
			r.Post("/", h.Create())
			r.Put("/{carbon_pk}", h.Update())
			r.Delete("/{carbon_pk}", h.Delete())
		})
	})
}
