package storage

import (
	"errors"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
)

const Gauge = "gauge"
const Counter = "counter"

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() Repository {
	return &Storage{
		&MemStorage{
			gauge:   map[string]float64{},
			counter: map[string]int64{},
		},
	}
}

type Storage struct {
	*MemStorage
}

type Repository interface {
	AddGauge(m metrics.Data) error
	AddCounter(m metrics.Data) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAllMetrics() *MemStorage
}

func (s *Storage) AddGauge(m metrics.Data) error {
	if m.Type == Gauge {
		s.gauge[m.Name] = m.Value
	} else {
		return errors.New("invalid metric type")
	}
	return nil
}

func (s *Storage) AddCounter(m metrics.Data) error {
	if m.Type == Counter {
		s.counter[m.Name] += int64(m.Value)
	} else {
		return errors.New("invalid metric type")
	}
	return nil
}

func (s *Storage) GetGauge(name string) (float64, error) {
	v, ok := s.gauge[name]
	if !ok {
		return 0, fmt.Errorf("metric %s not found", name)
	}
	return v, nil
}

func (s *Storage) GetCounter(name string) (int64, error) {
	v, ok := s.counter[name]
	if !ok {
		return 0, fmt.Errorf("metric %s not found", name)
	}
	return v, nil
}

func (s *Storage) GetAllMetrics() *MemStorage {
	return s.MemStorage
}
