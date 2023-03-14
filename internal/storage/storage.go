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

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
}

func (ms *MemStorage) AddGauge(m metrics.Data) error {
	if m.Type == Gauge {
		ms.gauge[m.Name] = m.Value
	} else {
		return errors.New("invalid metric type")
	}
	return nil
}

func (ms *MemStorage) AddCounter(m metrics.Data) error {
	if m.Type == Counter {
		ms.counter[m.Name] += int64(m.Value)
	} else {
		return errors.New("invalid metric type")
	}
	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	v, ok := ms.gauge[name]
	if !ok {
		return 0, fmt.Errorf("metric %s not found", name)
	}
	return v, nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	v, ok := ms.counter[name]
	if !ok {
		return 0, fmt.Errorf("metric %s not found", name)
	}
	return v, nil
}

func (ms *MemStorage) GetAllMetrics() ([]byte, error) {
	m := &MemStorage{
		map[string]float64{},
		map[string]int64{},
	}
	for k, v := range ms.gauge {
		m.gauge[k] = v
	}
	for k, v := range ms.counter {
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
