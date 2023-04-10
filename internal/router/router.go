package router

import (
	"github.com/eugeniylennik/alertics/internal/handlers"
	mw "github.com/eugeniylennik/alertics/internal/middleware"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(store *storage.MemStorage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.DefaultLogger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Use(mw.ContentJson)
	r.Use(mw.CompressGzip)
	r.Use(mw.DecompressGzip)

	r.Get("/", handlers.GetMetrics(store))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.RecordMetricsByJSON(store))
		r.Post("/{type}/{name}/{value}", handlers.RecordMetrics(store))
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetSpecificMetricJSON(store))
		r.Get("/{type}/{name}", handlers.GetSpecificMetric(store))
	})
	return r
}
