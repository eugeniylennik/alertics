package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eugeniylennik/alertics/internal/metrics"
	"github.com/eugeniylennik/alertics/internal/server"
	"github.com/eugeniylennik/alertics/internal/storage/file"
	"log"
	"sync"
)

const Gauge = "gauge"
const Counter = "counter"

type MemStorage struct {
	mux           sync.Mutex
	gauge         map[string]float64
	counter       map[string]int64
	isStoreToFile bool
	writer        *file.Writer
}

func NewMemStorage(cfg *server.Server) *MemStorage {
	w, err := file.NewWriter(cfg.StoreFile)
	if err != nil {
		log.Println(err)
	}
	return &MemStorage{
		gauge:         map[string]float64{},
		counter:       map[string]int64{},
		isStoreToFile: cfg.StoreInterval == 0,
		writer:        w,
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
	if ms.isStoreToFile {
		b, _ := json.Marshal(m)
		if err := ms.writer.WriteMetrics(b); err != nil {
			log.Println(err)
		}
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
	if ms.isStoreToFile {
		b, _ := json.Marshal(m)
		if err := ms.writer.WriteMetrics(b); err != nil {
			log.Println(err)
		}
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
