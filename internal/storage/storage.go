package storage

import (
	"errors"
	"github.com/eugeniylennik/alertics/internal/metrics"
)

const Gauge = "gauge"
const Counter = "counter"

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Repository interface {
	AddMetrics(m metrics.Data) error
}

func (ms *MemStorage) AddMetrics(m metrics.Data) error {
	switch m.Type {
	case Gauge:
		ms.Gauge[m.Name] = m.Value
	case Counter:
		ms.Counter[m.Name] += int64(m.Value)
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func NewRepository() *MemStorage {
	return &MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}
