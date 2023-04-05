package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Repository interface {
	AddGauge(m metrics.Data) error
	AddCounter(m metrics.Data) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAllMetrics() ([]byte, error)
}

func MiddlewareJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			if _, err := io.WriteString(w, err.Error()); err != nil {
				return
			}
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
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

func RecordMetricsByJSON(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m metrics.Metrics
		var d metrics.Data

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		d = metrics.Data{
			Name: m.ID,
			Type: m.MType,
		}

		if m.MType == storage.Gauge {
			d.Value = *m.Value
		} else {
			d.Value = float64(*m.Delta)
		}

		switch d.Type {
		case storage.Gauge:
			_ = repo.AddGauge(d)
			*m.Value, _ = repo.GetGauge(d.Name)
		case storage.Counter:
			_ = repo.AddCounter(d)
			*m.Delta, _ = repo.GetCounter(d.Name)
		}

		result, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func GetSpecificMetricJSON(repo Repository) http.HandlerFunc {
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

		switch m.MType {
		case storage.Gauge:
			v, err := repo.GetGauge(m.ID)
			if err != nil {
				fmt.Println("ERROR", err)
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
				fmt.Println(string(b))
				fmt.Println("ERROR: ", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		case storage.Counter:
			v, err := repo.GetCounter(m.ID)
			if err != nil {
				fmt.Println("ERROR: ", err)
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
				fmt.Println(string(b))
				fmt.Println("ERROR: ", err)
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
