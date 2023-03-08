package handlers

import (
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type Storage struct {
	storage.MemStorage
}

func (s *Storage) RecordMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path[len("/update/"):]
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	v, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := metrics.Data{
		Type:  parts[0],
		Name:  parts[1],
		Value: v,
	}

	if m.Type != "gauge" || m.Type != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err := s.Record(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func NewStorage() Storage {
	return Storage{
		storage.NewRepository(),
	}
}
