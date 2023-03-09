package router

import (
	"github.com/eugeniylennik/alertics/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	h := handlers.NewStorage()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.RecordMetrics)
	})
	r.Get("/", h.GetMetrics)
	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", h.GetSpecificMetric)
	})
	return r
}
