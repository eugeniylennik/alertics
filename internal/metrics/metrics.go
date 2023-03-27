package metrics

import (
	"encoding/json"
	"runtime"
)

var (
	inc int64 = 0
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"MType"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type ListMetrics []Metrics

func (m ListMetrics) MarshalJSON() ([]byte, error) {
	result := make([]Metrics, 0)
	for _, v := range m {
		result = append(result, v)
	}
	return json.Marshal(&result)
}

func CollectMetrics() ListMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	alloc := float64(memStats.Alloc)
	buckHasSys := float64(memStats.BuckHashSys)
	frees := float64(memStats.Frees)
	gccpuFraction := memStats.GCCPUFraction
	gcsys := float64(memStats.GCSys)
	heapAlloc := float64(memStats.HeapAlloc)
	heapIdle := float64(memStats.HeapIdle)
	heapInuse := float64(memStats.HeapInuse)
	heapObjects := float64(memStats.HeapObjects)
	heapReleased := float64(memStats.HeapReleased)
	heapSys := float64(memStats.HeapSys)
	lastGC := float64(memStats.LastGC)
	lookups := float64(memStats.Lookups)
	mCacheInuse := float64(memStats.MCacheSys)
	mCacheSys := float64(memStats.MCacheSys)
	mallocs := float64(memStats.Mallocs)
	nextGC := float64(memStats.NextGC)
	numForcedGC := float64(memStats.NumForcedGC)
	numGC := float64(memStats.NumGC)
	otherSys := float64(memStats.OtherSys)
	pauseTotalNs := float64(memStats.PauseTotalNs)
	stackInuse := float64(memStats.StackInuse)
	stackSys := float64(memStats.StackSys)
	sys := float64(memStats.Sys)
	totalAlloc := float64(memStats.TotalAlloc)
	randomValue := float64(memStats.TotalAlloc)

	inc++
	// Create a slice of metrics.
	metrics := []Metrics{
		{ID: "Alloc", MType: "gauge", Value: &alloc},
		{ID: "BuckHashSys", MType: "gauge", Value: &buckHasSys},
		{ID: "Frees", MType: "gauge", Value: &frees},
		{ID: "GCCPUFraction", MType: "gauge", Value: &gccpuFraction},
		{ID: "GCSys", MType: "gauge", Value: &gcsys},
		{ID: "HeapAlloc", MType: "gauge", Value: &heapAlloc},
		{ID: "HeapIdle", MType: "gauge", Value: &heapIdle},
		{ID: "HeapInuse", MType: "gauge", Value: &heapInuse},
		{ID: "HeapObjects", MType: "gauge", Value: &heapObjects},
		{ID: "HeapReleased", MType: "gauge", Value: &heapReleased},
		{ID: "HeapSys", MType: "gauge", Value: &heapSys},
		{ID: "LastGC", MType: "gauge", Value: &lastGC},
		{ID: "Lookups", MType: "gauge", Value: &lookups},
		{ID: "MCacheInuse", MType: "gauge", Value: &mCacheInuse},
		{ID: "MCacheSys", MType: "gauge", Value: &mCacheSys},
		{ID: "Mallocs", MType: "gauge", Value: &mallocs},
		{ID: "NextGC", MType: "gauge", Value: &nextGC},
		{ID: "NumForcedGC", MType: "gauge", Value: &numForcedGC},
		{ID: "NumGC", MType: "gauge", Value: &numGC},
		{ID: "OtherSys", MType: "gauge", Value: &otherSys},
		{ID: "PauseTotalNs", MType: "gauge", Value: &pauseTotalNs},
		{ID: "StackInuse", MType: "gauge", Value: &stackInuse},
		{ID: "StackSys", MType: "gauge", Value: &stackSys},
		{ID: "Sys", MType: "gauge", Value: &sys},
		{ID: "TotalAlloc", MType: "gauge", Value: &totalAlloc},
		{ID: "PollCount", MType: "counter", Delta: &inc},
		{ID: "RandomValue", MType: "gauge", Value: &randomValue},
	}

	return metrics
}
