package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
)

func (w MetricsWorker) collectMetrics(_ context.Context, m *entity.Metrics) error {
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
	m.RandomValue = entity.Gauge(rand.Float64())

	return nil
}

func (w MetricsWorker) sendMetrics(ctx context.Context, m *entity.Metrics) error {
	client := &http.Client{}

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
		if err := sendMetric(ctx, sendMetricOpts{
			client:        client,
			serverAddress: w.Opts.ServerAddr,
			metric:        metric,
		}); err != nil {
			return err
		}
	}
	for name, value := range counterMetrics {
		metric := entity.MetricDTO{
			Type:  enum.CounterMetricType,
			Name:  name,
			Value: fmt.Sprintf("%v", value),
		}
		if err := sendMetric(ctx, sendMetricOpts{
			client:        client,
			serverAddress: w.Opts.ServerAddr,
			metric:        metric,
		}); err != nil {
			return err
		}
	}

	return nil
}

type sendMetricOpts struct {
	client        *http.Client
	serverAddress string
	metric        entity.MetricDTO
}

func sendMetric(ctx context.Context, opts sendMetricOpts) error {
	metric := opts.metric
	url := fmt.Sprintf("%s/update/%s/%s/%s", opts.serverAddress, metric.Type, metric.Name, metric.Value)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := opts.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if cErr := Body.Close(); cErr != nil {
			log.Println(cErr)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %s\n", resp.Status)
	}

	return nil
}
