package handlers

import (
	"encoding/json"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Storage struct {
	storage.Repository
}

func (s *Storage) RecordMetrics(w http.ResponseWriter, r *http.Request) {
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
		_ = s.AddGauge(m)
	case storage.Counter:
		_ = s.AddCounter(m)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (s *Storage) GetSpecificMetric(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	switch typeMetric {
	case storage.Gauge:
		v, err := s.GetGauge(name)
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
		v, err := s.GetCounter(name)
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

func (s *Storage) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := s.GetAllMetrics()
	b, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func NewHandler(r storage.Repository) *Storage {
	return &Storage{
		r,
	}
}
