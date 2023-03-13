package storage

import (
	"encoding/json"
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
	GetAllMetrics() ([]byte, error)
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

func (s *Storage) GetAllMetrics() ([]byte, error) {
	m := &MemStorage{
		map[string]float64{},
		map[string]int64{},
	}
	for k, v := range s.MemStorage.gauge {
		m.gauge[k] = v
	}
	for k, v := range s.MemStorage.counter {
		m.counter[k] = v
	}

	b, err := json.Marshal(struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}{
		Gauge:   m.gauge,
		Counter: m.counter,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}
