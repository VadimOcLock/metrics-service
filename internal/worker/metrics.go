package worker

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/hashutil"

	"github.com/VadimOcLock/metrics-service/internal/compress"

	"github.com/go-resty/resty/v2"

	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
)

const updateAPIEndpoint = "/updates/"

func (w *MetricsWorker) collectRuntimeMetrics(_ context.Context) (entity.MetricsData, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var m entity.MetricsData

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
		return entity.MetricsData{}, fmt.Errorf("worker.collectRuntimeMetrics: %w", err)
	}
	bigFloat := new(big.Float).Quo(new(big.Float).SetInt(randomInt), big.NewFloat(10000))
	randVal, _ := bigFloat.Float64()
	m.RandomValue = entity.Gauge(randVal)

	return m, nil
}

func (w *MetricsWorker) collectRuntimeMetricsLoop(ctx context.Context, errCh chan error) {
	ticker := time.NewTicker(w.Opts.PoolInterval)
	defer ticker.Stop()
	for t := range ticker.C {
		log.Debug().Msgf("runtime metrics collector start at: %s", t.String())
		data, err := w.collectRuntimeMetrics(ctx)
		if err != nil {
			errCh <- fmt.Errorf("collect runtime metrics err: %w", err)
			continue
		}
		select {
		case <-ctx.Done():
			log.Debug().Msg("runtime metrics collector finished")
			return
		case w.MetricsCh <- data:
		default:
		}
	}
}

func (w *MetricsWorker) collectSystemMetrics(_ context.Context) (entity.MetricsData, error) {
	var metrics entity.MetricsData

	v, err := mem.VirtualMemory()
	if err != nil {
		return metrics, err
	}
	metrics.TotalMemory = entity.Gauge(v.Total)
	metrics.FreeMemory = entity.Gauge(v.Free)

	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		return metrics, err
	}
	if len(cpuUtilization) > 0 {
		metrics.CPUUtilization1 = entity.Gauge(cpuUtilization[0])
	}

	return metrics, nil
}

func (w *MetricsWorker) collectSystemMetricsLoop(ctx context.Context, errCh chan error) {
	ticker := time.NewTicker(w.Opts.PoolInterval)
	defer ticker.Stop()
	for t := range ticker.C {
		log.Debug().Msgf("system metrics collector start at: %s", t.String())
		data, err := w.collectSystemMetrics(ctx)
		if err != nil {
			errCh <- fmt.Errorf("collect system metrics err: %w", err)
			continue
		}
		select {
		case <-ctx.Done():
			log.Debug().Msg("system metrics collector finished")
			return
		case w.MetricsCh <- data:
		default:
		}
	}
}

func (w *MetricsWorker) sendMetrics(ctx context.Context, m entity.MetricsData) error {
	gauges, counters := buildMetricsMap(m)
	metricsBatch, err := buildMetricsBatch(gauges, counters)
	if err != nil {
		return err
	}
	return sendMetricRequest(ctx, sendMetricRequestOpts{
		ServerAddress:      w.Opts.ServerAddr,
		SecretSignatureKey: w.Opts.SecretSignatureKey,
		batch:              metricsBatch,
	})
}

func buildMetricsBatch(gs map[string]entity.Gauge, cs map[string]entity.Counter) ([]entity.Metrics, error) {
	outLen := len(gs) + len(cs)
	if outLen == 0 {
		return nil, errorz.ErrTrySendEmptyData
	}
	res := make([]entity.Metrics, outLen)
	for name, val := range gs {
		vl := float64(val)
		if val == 0 {
			continue
		}
		m := entity.Metrics{
			ID:    name,
			MType: enum.GaugeMetricType,
			Value: &vl,
		}
		res = append(res, m)
	}
	for name, val := range cs {
		vl := int64(val)
		if val == 0 {
			continue
		}
		m := entity.Metrics{
			ID:    name,
			MType: enum.CounterMetricType,
			Delta: &vl,
		}
		res = append(res, m)
	}

	return res, nil
}

func buildMetricsMap(m entity.MetricsData) (map[string]entity.Gauge, map[string]entity.Counter) {
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

		enum.TotalMemoryName:     m.TotalMemory,
		enum.FreeMemoryName:      m.FreeMemory,
		enum.CPUUtilization1Name: m.CPUUtilization1,
	}
	counterMetrics := map[string]entity.Counter{
		enum.PollCountMetricName: m.PollCount,
	}

	return gaugeMetrics, counterMetrics
}

type sendMetricRequestOpts struct {
	ServerAddress      string
	SecretSignatureKey string
	batch              []entity.Metrics
}

func sendMetricRequest(ctx context.Context, opts sendMetricRequestOpts) error {
	client := resty.New()
	url := opts.ServerAddress + updateAPIEndpoint

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(opts.batch); err != nil {
		return fmt.Errorf("worker.sendMetricReq: %w", err)
	}

	body, err := compress.GZipCompress(buf.Bytes())
	if err != nil {
		return fmt.Errorf("worker.sendMetricReq: %w", err)
	}

	req := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip")

	if opts.SecretSignatureKey != "" {
		hash := hashutil.ComputeHMAC(buf.Bytes(), opts.SecretSignatureKey)
		req.SetHeader("HashSHA256", hash)
	}

	resp, err := req.
		SetBody(body).
		Post(url)
	if err != nil {
		return fmt.Errorf("worker.sendMetricReq: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error().Msg(string(resp.Body()))

		return errorz.ErrSendMetricStatusNotOK
	}

	return nil
}

func (w *MetricsWorker) sendMetricsLoop(ctx context.Context, errCh chan error) {
	ticker := time.NewTicker(w.Opts.ReportInterval)
	defer ticker.Stop()
	for t := range ticker.C {
		log.Debug().Msgf("metrics sender start at: %s", t.String())
		select {
		case <-ctx.Done():
			log.Debug().Msg("metrics sender worker finished")
			return
		case data, ok := <-w.MetricsCh:
			if ok {
				errCh <- w.sendMetrics(ctx, data)
			}
		default:
		}
	}
}
