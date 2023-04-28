package router

import (
	"github.com/eugeniylennik/alertics/internal/handlers"
	mw "github.com/eugeniylennik/alertics/internal/middleware"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/eugeniylennik/alertics/internal/storage/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(store *storage.MemStorage, db *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.DefaultLogger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Use(mw.ContentTypeJSON)
	r.Use(mw.CompressGzip)
	r.Use(mw.DecompressGzip)

	dbStore := database.NewStorage(db)

	r.Get("/", handlers.GetMetrics(store))
	r.Get("/ping", handlers.HealthCheckDB(db))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.RecordMetricsByJSON(store, dbStore))
		r.Post("/{type}/{name}/{value}", handlers.RecordMetrics(store))
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.GetSpecificMetricJSON(store))
		r.Get("/{type}/{name}", handlers.GetSpecificMetric(store))
	})
	return r
}
