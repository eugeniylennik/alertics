package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/server"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/eugeniylennik/alertics/internal/storage/database"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"strconv"
)

var cfg = server.InitConfigServer()

type Repository interface {
	AddGauge(m metrics.Data) error
	AddCounter(m metrics.Data) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAllMetrics() ([]byte, error)
}

func RecordMetrics(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeMetric := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		if typeMetric != storage.Gauge && typeMetric != storage.Counter {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m := metrics.Data{
			Type:  typeMetric,
			Name:  name,
			Value: v,
		}
		switch typeMetric {
		case storage.Gauge:
			_ = repo.AddGauge(m)
		case storage.Counter:
			_ = repo.AddCounter(m)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func RecordMetricsByJSON(repo Repository, db database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metrics.Metrics
		var d metrics.Data

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if m.Hash == "" {
			http.Error(w, fmt.Sprintf("hash empty"), http.StatusBadRequest)
			return
		}

		_, err := m.IsHashesEquals()
		//if !ok {
		//	http.Error(w, fmt.Sprintf("hashes is not equals"), http.StatusBadRequest)
		//	return
		//}

		d = metrics.Data{
			Name: m.ID,
			Type: m.MType,
		}

		if m.MType == storage.Gauge {
			d.Value = *m.Value
		} else {
			d.Value = float64(*m.Delta)
		}

		if cfg.Dsn == "" {
			switch d.Type {
			case storage.Gauge:
				_ = repo.AddGauge(d)
				*m.Value, _ = repo.GetGauge(d.Name)
			case storage.Counter:
				_ = repo.AddCounter(d)
				*m.Delta, _ = repo.GetCounter(d.Name)
			}
		} else {
			err = db.InsertMetrics(context.Background(), m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		result, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func GetSpecificMetricJSON(repo Repository, db database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var m metrics.Metrics
		if err := json.Unmarshal(b, &m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if cfg.Dsn == "" {
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
				b, err := json.Marshal(r)
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
				b, err := json.Marshal(r)
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
		} else {
			r, err := db.SelectMetricById(context.Background(), m)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			b, err := json.Marshal(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}
}

func GetSpecificMetric(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeMetric := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch typeMetric {
		case storage.Gauge:
			v, err := repo.GetGauge(name)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			b, err := json.Marshal(v)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		case storage.Counter:
			v, err := repo.GetCounter(name)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			b, err := json.Marshal(v)
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(m)
	}
}

func HealthCheckDB(pgx *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pgx.Ping(context.TODO())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func RecordMetricsBatch(db database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m []metrics.Metrics

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err := db.InsertMetricsStatement(context.Background(), m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		result, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
