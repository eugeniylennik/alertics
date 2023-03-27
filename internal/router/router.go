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
	r.Use(middleware.StripSlashes)
	r.Use(handlers.MiddlewareJson)

	r.Get("/", handlers.GetMetrics(m))
	r.Post("/update", handlers.RecordMetrics(m))
	r.Post("/value", handlers.GetSpecificMetric(m))
	return r
}
