package entity

import (
	"strconv"

	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
)

type (
	Gauge   float64
	Counter int64
)

type MetricsData struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	PollCount     Counter
	RandomValue   Gauge
}

type MetricDTO struct {
	Type  string
	Name  string
	Value string
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) Valid() error {
	if m.MType == enum.GaugeMetricType || m.MType == enum.CounterMetricType {
		return nil
	}

	return errorz.ErrUndefinedMetricType
}

func (m *Metrics) MetricValue() string {
	if m.MType == "gauge" && m.Value != nil {
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	}
	if m.MType == "counter" && m.Delta != nil {
		return strconv.FormatInt(*m.Delta, 10)
	}

	return ""
}
