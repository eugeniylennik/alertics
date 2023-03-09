package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Storage struct {
	m *storage.MemStorage
}

type ApiResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *Storage) RecordMetrics(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if typeMetric != "gauge" && typeMetric != "counter" {
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

	if err := s.m.AddMetrics(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (s *Storage) GetSpecificMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	typeMetric := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	var value float64
	switch typeMetric {
	case storage.Gauge:
		v, ok := s.m.Gauge[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			b, _ := json.Marshal(
				ApiResponse{
					Message: fmt.Sprintf("Name %s is not found", name),
				})
			w.Write(b)
			return
		}
		value = v
	case storage.Counter:
		v, ok := s.m.Counter[name]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			b, _ := json.Marshal(
				ApiResponse{
					Message: fmt.Sprintf("Name %s is not found", name),
				})
			w.Write(b)
			return
		}
		value = float64(v)
	default:
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(ApiResponse{Message: fmt.Sprintf("Type %s is not found", typeMetric)})
		w.Write(b)
		return
	}

	result := struct {
		Name  string
		Value float64
	}{
		Name:  name,
		Value: value,
	}

	b, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Storage) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(s.m)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func NewStorage() Storage {
	return Storage{
		storage.NewRepository(),
	}
}
