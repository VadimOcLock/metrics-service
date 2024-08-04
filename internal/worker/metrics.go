package worker

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"runtime"

	"github.com/go-resty/resty/v2"

	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
)

func (w *MetricsWorker) collectMetrics(_ context.Context, m *entity.MetricsData) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.Alloc = entity.Gauge(memStats.Alloc)
	m.BuckHashSys = entity.Gauge(memStats.BuckHashSys)
	m.Frees = entity.Gauge(memStats.Frees)
	m.GCCPUFraction = entity.Gauge(memStats.GCCPUFraction)
	m.GCSys = entity.Gauge(memStats.GCSys)
	m.HeapAlloc = entity.Gauge(memStats.HeapAlloc)
	m.HeapIdle = entity.Gauge(memStats.HeapIdle)
	m.HeapInuse = entity.Gauge(memStats.HeapInuse)
	m.HeapObjects = entity.Gauge(memStats.HeapObjects)
	m.HeapReleased = entity.Gauge(memStats.HeapReleased)
	m.HeapSys = entity.Gauge(memStats.HeapSys)
	m.LastGC = entity.Gauge(memStats.LastGC)
	m.Lookups = entity.Gauge(memStats.Lookups)
	m.MCacheInuse = entity.Gauge(memStats.MCacheInuse)
	m.MCacheSys = entity.Gauge(memStats.MCacheSys)
	m.MSpanInuse = entity.Gauge(memStats.MSpanInuse)
	m.MSpanSys = entity.Gauge(memStats.MSpanSys)
	m.Mallocs = entity.Gauge(memStats.Mallocs)
	m.NextGC = entity.Gauge(memStats.NextGC)
	m.NumForcedGC = entity.Gauge(memStats.NumForcedGC)
	m.NumGC = entity.Gauge(memStats.NumGC)
	m.OtherSys = entity.Gauge(memStats.OtherSys)
	m.PauseTotalNs = entity.Gauge(memStats.PauseTotalNs)
	m.StackInuse = entity.Gauge(memStats.StackInuse)
	m.StackSys = entity.Gauge(memStats.StackSys)
	m.Sys = entity.Gauge(memStats.Sys)
	m.TotalAlloc = entity.Gauge(memStats.TotalAlloc)

	m.PollCount++

	maxInt := big.NewInt(1000000)
	randomInt, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		return fmt.Errorf("worker.collectMetrics: %w", err)
	}
	bigFloat := new(big.Float).Quo(new(big.Float).SetInt(randomInt), big.NewFloat(10000))
	randVal, _ := bigFloat.Float64()
	m.RandomValue = entity.Gauge(randVal)

	return nil
}

func (w *MetricsWorker) sendMetrics(ctx context.Context, m *entity.MetricsData) error {
	client := resty.New()

	gaugeMetrics := map[string]entity.Gauge{
		enum.AllocMetricName:         m.Alloc,
		enum.BuckHashSysMetricName:   m.BuckHashSys,
		enum.FreesMetricName:         m.Frees,
		enum.GCCPUFractionMetricName: m.GCCPUFraction,
		enum.GCSysMetricName:         m.GCSys,
		enum.HeapAllocMetricName:     m.HeapAlloc,
		enum.HeapIdleMetricName:      m.HeapIdle,
		enum.HeapInuseMetricName:     m.HeapInuse,
		enum.HeapObjectsMetricName:   m.HeapObjects,
		enum.HeapReleasedMetricName:  m.HeapReleased,
		enum.HeapSysMetricName:       m.HeapSys,
		enum.LastGCMetricName:        m.LastGC,
		enum.LookupsMetricName:       m.Lookups,
		enum.MCacheInuseMetricName:   m.MCacheInuse,
		enum.MCacheSysMetricName:     m.MCacheSys,
		enum.MSpanInuseMetricName:    m.MSpanInuse,
		enum.MSpanSysMetricName:      m.MSpanSys,
		enum.MallocsMetricName:       m.Mallocs,
		enum.NextGCMetricName:        m.NextGC,
		enum.NumForcedGCMetricName:   m.NumForcedGC,
		enum.NumGCMetricName:         m.NumGC,
		enum.OtherSysMetricName:      m.OtherSys,
		enum.PauseTotalNsMetricName:  m.PauseTotalNs,
		enum.StackInuseMetricName:    m.StackInuse,
		enum.StackSysMetricName:      m.StackSys,
		enum.SysMetricName:           m.Sys,
		enum.TotalAllocMetricName:    m.TotalAlloc,
		enum.RandomValueMetricName:   m.RandomValue,
	}
	counterMetrics := map[string]entity.Counter{
		enum.PollCountMetricName: m.PollCount,
	}

	for name, value := range gaugeMetrics {
		metric := entity.MetricDTO{
			Type:  enum.GaugeMetricType,
			Name:  name,
			Value: fmt.Sprintf("%v", value),
		}
		if err := SendMetric(ctx, SendMetricOpts{
			Client:        client,
			ServerAddress: w.Opts.ServerAddr,
			Metric:        metric,
		}); err != nil {
			return fmt.Errorf("worker.sendMetrics: %w", err)
		}
	}
	for name, value := range counterMetrics {
		metric := entity.MetricDTO{
			Type:  enum.CounterMetricType,
			Name:  name,
			Value: fmt.Sprintf("%v", value),
		}
		if err := SendMetric(ctx, SendMetricOpts{
			Client:        client,
			ServerAddress: w.Opts.ServerAddr,
			Metric:        metric,
		}); err != nil {
			return fmt.Errorf("worker.sendMetrics: %w", err)
		}
	}

	return nil
}

type SendMetricOpts struct {
	Client        *resty.Client
	ServerAddress string
	Metric        entity.MetricDTO
}

func SendMetric(_ context.Context, opts SendMetricOpts) error {
	metric := opts.Metric
	url := fmt.Sprintf("%s/update/%s/%s/%s", opts.ServerAddress, metric.Type, metric.Name, metric.Value)

	resp, err := opts.Client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		return fmt.Errorf("worker.SendMetric: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return errorz.ErrSendMetricStatusNotOK
	}

	return nil
}

//package worker
//
//import (
//	"context"
//	"crypto/rand"
//	"encoding/json"
//	"fmt"
//	"github.com/rs/zerolog/log"
//	"math/big"
//	"net/http"
//	"runtime"
//
//	"github.com/go-resty/resty/v2"
//
//	"github.com/VadimOcLock/metrics-service/internal/errorz"
//
//	"github.com/VadimOcLock/metrics-service/internal/entity"
//	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
//)
//
//func (w *MetricsWorker) collectMetrics(_ context.Context, m *entity.MetricsData) error {
//	var memStats runtime.MemStats
//	runtime.ReadMemStats(&memStats)
//
//	m.Alloc = entity.Gauge(memStats.Alloc)
//	m.BuckHashSys = entity.Gauge(memStats.BuckHashSys)
//	m.Frees = entity.Gauge(memStats.Frees)
//	m.GCCPUFraction = entity.Gauge(memStats.GCCPUFraction)
//	m.GCSys = entity.Gauge(memStats.GCSys)
//	m.HeapAlloc = entity.Gauge(memStats.HeapAlloc)
//	m.HeapIdle = entity.Gauge(memStats.HeapIdle)
//	m.HeapInuse = entity.Gauge(memStats.HeapInuse)
//	m.HeapObjects = entity.Gauge(memStats.HeapObjects)
//	m.HeapReleased = entity.Gauge(memStats.HeapReleased)
//	m.HeapSys = entity.Gauge(memStats.HeapSys)
//	m.LastGC = entity.Gauge(memStats.LastGC)
//	m.Lookups = entity.Gauge(memStats.Lookups)
//	m.MCacheInuse = entity.Gauge(memStats.MCacheInuse)
//	m.MCacheSys = entity.Gauge(memStats.MCacheSys)
//	m.MSpanInuse = entity.Gauge(memStats.MSpanInuse)
//	m.MSpanSys = entity.Gauge(memStats.MSpanSys)
//	m.Mallocs = entity.Gauge(memStats.Mallocs)
//	m.NextGC = entity.Gauge(memStats.NextGC)
//	m.NumForcedGC = entity.Gauge(memStats.NumForcedGC)
//	m.NumGC = entity.Gauge(memStats.NumGC)
//	m.OtherSys = entity.Gauge(memStats.OtherSys)
//	m.PauseTotalNs = entity.Gauge(memStats.PauseTotalNs)
//	m.StackInuse = entity.Gauge(memStats.StackInuse)
//	m.StackSys = entity.Gauge(memStats.StackSys)
//	m.Sys = entity.Gauge(memStats.Sys)
//	m.TotalAlloc = entity.Gauge(memStats.TotalAlloc)
//
//	m.PollCount++
//
//	maxInt := big.NewInt(1000000)
//	randomInt, err := rand.Int(rand.Reader, maxInt)
//	if err != nil {
//		return fmt.Errorf("worker.collectMetrics: %w", err)
//	}
//	bigFloat := new(big.Float).Quo(new(big.Float).SetInt(randomInt), big.NewFloat(10000))
//	randVal, _ := bigFloat.Float64()
//	m.RandomValue = entity.Gauge(randVal)
//
//	return nil
//}
//
//func (w *MetricsWorker) sendMetrics(ctx context.Context, m *entity.MetricsData) error {
//	client := resty.New()
//
//	gaugeMetrics := map[string]entity.Gauge{
//		enum.AllocMetricName:         m.Alloc,
//		enum.BuckHashSysMetricName:   m.BuckHashSys,
//		enum.FreesMetricName:         m.Frees,
//		enum.GCCPUFractionMetricName: m.GCCPUFraction,
//		enum.GCSysMetricName:         m.GCSys,
//		enum.HeapAllocMetricName:     m.HeapAlloc,
//		enum.HeapIdleMetricName:      m.HeapIdle,
//		enum.HeapInuseMetricName:     m.HeapInuse,
//		enum.HeapObjectsMetricName:   m.HeapObjects,
//		enum.HeapReleasedMetricName:  m.HeapReleased,
//		enum.HeapSysMetricName:       m.HeapSys,
//		enum.LastGCMetricName:        m.LastGC,
//		enum.LookupsMetricName:       m.Lookups,
//		enum.MCacheInuseMetricName:   m.MCacheInuse,
//		enum.MCacheSysMetricName:     m.MCacheSys,
//		enum.MSpanInuseMetricName:    m.MSpanInuse,
//		enum.MSpanSysMetricName:      m.MSpanSys,
//		enum.MallocsMetricName:       m.Mallocs,
//		enum.NextGCMetricName:        m.NextGC,
//		enum.NumForcedGCMetricName:   m.NumForcedGC,
//		enum.NumGCMetricName:         m.NumGC,
//		enum.OtherSysMetricName:      m.OtherSys,
//		enum.PauseTotalNsMetricName:  m.PauseTotalNs,
//		enum.StackInuseMetricName:    m.StackInuse,
//		enum.StackSysMetricName:      m.StackSys,
//		enum.SysMetricName:           m.Sys,
//		enum.TotalAllocMetricName:    m.TotalAlloc,
//		enum.RandomValueMetricName:   m.RandomValue,
//	}
//	counterMetrics := map[string]entity.Counter{
//		enum.PollCountMetricName: m.PollCount,
//	}
//
//	for name, value := range gaugeMetrics {
//		metric := entity.Metrics{
//			ID:    name,
//			MType: enum.GaugeMetricType,
//			Value: (*float64)(&value),
//		}
//		if err := SendMetric(ctx, SendMetricOpts{
//			Client:        client,
//			ServerAddress: w.Opts.ServerAddr,
//			Metric:        metric,
//		}); err != nil {
//			return fmt.Errorf("worker.sendMetrics: %w", err)
//		}
//	}
//	for name, value := range counterMetrics {
//		metric := entity.Metrics{
//			ID:    name,
//			MType: enum.CounterMetricType,
//			Delta: (*int64)(&value),
//		}
//		if err := SendMetric(ctx, SendMetricOpts{
//			Client:        client,
//			ServerAddress: w.Opts.ServerAddr,
//			Metric:        metric,
//		}); err != nil {
//			return fmt.Errorf("worker.sendMetrics: %w", err)
//		}
//	}
//
//	return nil
//}
//
//type SendMetricOpts struct {
//	Client        *resty.Client
//	ServerAddress string
//	Metric        entity.Metrics
//}
//
//func SendMetric(_ context.Context, opts SendMetricOpts) error {
//	url := opts.ServerAddress + "/update/"
//	//url := fmt.Sprintf("%s/update/%s/%s/%s", opts.ServerAddress,
//	//	opts.Metric.MType, opts.Metric.ID, opts.Metric.MetricValue())
//
//	body, err := json.Marshal(opts.Metric)
//	if err != nil {
//		return fmt.Errorf("worker.SendMetric: %w", err)
//	}
//
//	resp, err := opts.Client.R().
//		SetHeader("Content-Type", "application/json").
//		SetBody(body).
//		Post(url)
//	if err != nil {
//		log.Debug().Msg(string(body))
//		return fmt.Errorf("worker.SendMetric: %w", err)
//	}
//
//	if resp.StatusCode() != http.StatusOK {
//		return errorz.ErrSendMetricStatusNotOK
//	}
//
//	return nil
//}
