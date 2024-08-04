package entity

import (
	"errors"
	"fmt"
	"strconv"
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

type Metric struct {
	Type  string
	Name  string
	Value any
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// BuildMetrics - функция для создания объекта Metrics из MetricDTO.
func BuildMetrics(dto MetricDTO) (Metrics, error) {
	metric := Metrics{
		ID:    dto.Name,
		MType: dto.Type,
	}

	switch dto.Type {
	case "gauge":
		value, err := strconv.ParseFloat(dto.Value, 64)
		if err != nil {
			return Metrics{}, fmt.Errorf("invalid gauge value: %v", err)
		}
		metric.Value = &value
	case "counter":
		delta, err := strconv.ParseInt(dto.Value, 10, 64)
		if err != nil {
			return Metrics{}, fmt.Errorf("invalid counter value: %v", err)
		}
		metric.Delta = &delta
	default:
		return Metrics{}, fmt.Errorf("unknown metric type: %s", dto.Type)
	}

	return metric, nil
}

// MetricValue - метод для получения значения метрики в виде строки.
func (m *Metrics) MetricValue() (string, error) {
	switch m.MType {
	case "gauge":
		if m.Value != nil {
			return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
		}

		return "", fmt.Errorf("value is nil for gauge type")
	case "counter":
		if m.Delta != nil {
			return strconv.FormatInt(*m.Delta, 10), nil
		}

		return "", fmt.Errorf("delta is nil for counter type")
	default:
		return "", fmt.Errorf("unknown metric type: %s", m.MType)
	}
}

// Valid - метод для проверки корректности метрики.
func (m *Metrics) Valid() error {
	switch m.MType {
	case "gauge":
		if m.Value == nil {
			return errors.New("value is nil for gauge type")
		}
	case "counter":
		if m.Delta == nil {
			return errors.New("delta is nil for counter type")
		}
	default:
		return fmt.Errorf("unknown metric type: %s", m.MType)
	}

	return nil
}
