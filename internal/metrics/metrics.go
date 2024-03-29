package metrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
)

var (
	inc int64 = 0
)

type Data struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func CollectMetrics() []Data {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	inc++
	// Create a slice of metrics.
	metrics := []Data{
		{Name: "Alloc", Type: "gauge", Value: float64(memStats.Alloc)},
		{Name: "BuckHashSys", Type: "gauge", Value: float64(memStats.BuckHashSys)},
		{Name: "Frees", Type: "gauge", Value: float64(memStats.Frees)},
		{Name: "GCCPUFraction", Type: "gauge", Value: memStats.GCCPUFraction},
		{Name: "GCSys", Type: "gauge", Value: float64(memStats.GCSys)},
		{Name: "HeapAlloc", Type: "gauge", Value: float64(memStats.HeapAlloc)},
		{Name: "HeapIdle", Type: "gauge", Value: float64(memStats.HeapIdle)},
		{Name: "HeapInuse", Type: "gauge", Value: float64(memStats.HeapInuse)},
		{Name: "HeapObjects", Type: "gauge", Value: float64(memStats.HeapObjects)},
		{Name: "HeapReleased", Type: "gauge", Value: float64(memStats.HeapReleased)},
		{Name: "HeapSys", Type: "gauge", Value: float64(memStats.HeapSys)},
		{Name: "LastGC", Type: "gauge", Value: float64(memStats.LastGC)},
		{Name: "Lookups", Type: "gauge", Value: float64(memStats.Lookups)},
		{Name: "MCacheInuse", Type: "gauge", Value: float64(memStats.MCacheInuse)},
		{Name: "MCacheSys", Type: "gauge", Value: float64(memStats.MCacheSys)},
		{Name: "Mallocs", Type: "gauge", Value: float64(memStats.Mallocs)},
		{Name: "NextGC", Type: "gauge", Value: float64(memStats.NextGC)},
		{Name: "NumForcedGC", Type: "gauge", Value: float64(memStats.NumForcedGC)},
		{Name: "NumGC", Type: "gauge", Value: float64(memStats.NumGC)},
		{Name: "OtherSys", Type: "gauge", Value: float64(memStats.OtherSys)},
		{Name: "PauseTotalNs", Type: "gauge", Value: float64(memStats.PauseTotalNs)},
		{Name: "StackInuse", Type: "gauge", Value: float64(memStats.StackInuse)},
		{Name: "StackSys", Type: "gauge", Value: float64(memStats.StackSys)},
		{Name: "Sys", Type: "gauge", Value: float64(memStats.Sys)},
		{Name: "TotalAlloc", Type: "gauge", Value: float64(memStats.TotalAlloc)},
		{Name: "PollCount", Type: "counter", Value: float64(inc)},
		{Name: "RandomValue", Type: "gauge", Value: float64(memStats.TotalAlloc)},
	}

	return metrics
}

func (m Metrics) IsHashesEquals() (bool, error) {
	h := hmac.New(sha256.New, []byte("key"))

	data, err := hex.DecodeString(m.Hash)
	if err != nil {
		return false, err
	}

	var msg string
	if m.MType == "counter" {
		msg = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	} else {
		msg = fmt.Sprintf("%s:gauge:%.f", m.ID, *m.Value)
	}

	h.Write([]byte(msg))
	sign := h.Sum(nil)

	if hmac.Equal(sign, data) {
		return true, nil
	} else {
		return false, nil
	}
}
