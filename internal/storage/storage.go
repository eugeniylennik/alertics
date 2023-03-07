package storage

import (
	"errors"
	"github.com/eugeniylennik/alertics/internal/metrics"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Repository interface {
	Record(m metrics.Data) error
}

func (ms *MemStorage) Record(m metrics.Data) error {
	switch m.Type {
	case "gauge":
		ms.Gauge[m.Name] = m.Value
	case "counter":
		ms.Counter[m.Name] += int64(m.Value)
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func NewRepository() MemStorage {
	return MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}
