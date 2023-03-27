package handlers

import (
	"encoding/json"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"net/http"
)

type Repository interface {
	AddGauge(m metrics.Metrics) error
	AddCounter(m metrics.Metrics) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAllMetrics() ([]byte, error)
}

func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RecordMetrics(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metrics.Metrics
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		switch m.MType {
		case storage.Gauge:
			_ = repo.AddGauge(m)
		case storage.Counter:
			_ = repo.AddCounter(m)
		}

		result, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func GetSpecificMetric(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metrics.Metrics
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		switch m.MType {
		case storage.Gauge:
			v, err := repo.GetGauge(m.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			r := metrics.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Value: &v,
			}
			b, err := json.MarshalIndent(r, "", " ")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		case storage.Counter:
			v, err := repo.GetCounter(m.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			r := metrics.Metrics{
				ID:    m.ID,
				MType: m.MType,
				Delta: &v,
			}
			b, err := json.MarshalIndent(r, "", " ")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func GetMetrics(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := repo.GetAllMetrics()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(m)
	}
}
