package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VadimOcLock/metrics-service/internal/config"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Listener struct {
	FileUpdater chan bool
	Cfg         *config.FileWriter
	store       metricservice.Store
}

func NewListener(ch chan bool, cfg *config.FileWriter, store metricservice.Store) *Listener {
	return &Listener{
		FileUpdater: ch,
		Cfg:         cfg,
		store:       store,
	}
}

func (l *Listener) saveMetrics(filePath string, metrics []entity.Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	var file *os.File
	if l.Cfg.Restore {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}(file)

	return nil
}

func (l *Listener) Run(ctx context.Context) error {
	if l.Cfg.StoreInterval == 0 {
		// Синхронно.
		for {
			select {
			case <-ctx.Done():
				err := l.writeMetrics(ctx)
				if err != nil {
					log.Error().Err(err)
				}
				return ctx.Err()
			case <-l.FileUpdater:
				err := l.writeMetrics(ctx)
				if err != nil {
					log.Error().Err(err)
				}
			}
		}
	} else {
		// Периодически.
		ticker := time.NewTicker(time.Duration(l.Cfg.StoreInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				err := l.writeMetrics(ctx)
				if err != nil {
					log.Error().Err(err)
				}
				return ctx.Err()
			case <-ticker.C:
				err := l.writeMetrics(ctx)
				if err != nil {
					log.Error().Err(err)
				}
			}
		}
	}
}

func (l *Listener) writeMetrics(ctx context.Context) error {
	data, err := l.metricsData(ctx)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	if err = l.saveMetrics(l.Cfg.FileStoragePath, data); err != nil {
		return err
	}

	return nil
}

func (l *Listener) metricsData(ctx context.Context) ([]entity.Metrics, error) {
	metricsData, err := l.store.FindAllMetrics(ctx, somestore.FindAllMetricsParams{})
	if err != nil {
		return nil, fmt.Errorf("file.metricsData: %w", err)
	}
	var metrics []entity.Metrics
	for _, md := range metricsData {
		m, err := entity.BuildMetrics(entity.MetricDTO{
			Type:  md.Type,
			Name:  md.Name,
			Value: fmt.Sprintf("%v", md.Value),
		})
		if err != nil {
			return nil, fmt.Errorf("file.metricsData: %w", err)
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
