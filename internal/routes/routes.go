package routes

import (
	"github.com/crafty-ezhik/carbonstats/internal/service_description"
	"github.com/crafty-ezhik/carbonstats/internal/statistics"
	"github.com/crafty-ezhik/carbonstats/internal/stats_data"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(r *chi.Mux, h service_description.ServiceDescriptionHandler, hs statistics.StatisticsHandler, sd stats_data.StatsDataHandler) {
	r.Route("/carbonstats", func(r chi.Router) {
		r.Route("/services_desc", func(r chi.Router) {
			r.Get("/{carbon_pk}", h.GetByCarbonPK())
			r.Get("/", h.GetAll())
			r.Post("/", h.Create())
			r.Put("/{carbon_pk}", h.Update())
			r.Delete("/{carbon_pk}", h.Delete())
		})

		r.Route("/stats", func(r chi.Router) {
			r.Get("/", hs.GetAll())
			r.Get("/get_stats", sd.GetStats())
			r.Post("/", hs.Create())
		})
	})
}
