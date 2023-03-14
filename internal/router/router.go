package router

import (
	"github.com/eugeniylennik/alertics/internal/handlers"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	m := storage.NewMemStorage()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", handlers.GetMetrics(m))
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handlers.RecordMetrics(m))
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handlers.GetSpecificMetric(m))
	})
	return r
}
