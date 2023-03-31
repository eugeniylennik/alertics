package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"log"
	"sync"
)

const Gauge = "gauge"
const Counter = "counter"

type MemStorage struct {
	mux     sync.Mutex
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
	ms.mux.Lock()
	defer ms.mux.Unlock()
	if m.Type == Gauge {
		ms.gauge[m.Name] = m.Value
	} else {
		return errors.New("invalid metric type")
	}
	return nil
}

func (ms *MemStorage) AddCounter(m metrics.Data) error {
	ms.mux.Lock()
	defer ms.mux.Unlock()
	if m.Type == Counter {
		ms.counter[m.Name] += int64(m.Value)
	} else {
		return errors.New("invalid metric type")
	}
	return nil
}

func (ms *MemStorage) GetGauge(name string) (float64, error) {
	log.Println("GETTING - ", name)
	v, ok := ms.gauge[name]
	log.Println("VALUE - ", v)
	if !ok {
		return 0, fmt.Errorf("metric %s not found", name)
	}
	return v, nil
}

func (ms *MemStorage) GetCounter(name string) (int64, error) {
	log.Println("GETTING - ", name)
	v, ok := ms.counter[name]
	log.Println("VALUE - ", v)
	if !ok {
		return 0, fmt.Errorf("metric %s not found", name)
	}
	return v, nil
}

func (ms *MemStorage) GetAllMetrics() ([]byte, error) {
	m := &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
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
