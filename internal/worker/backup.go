package worker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/rs/zerolog/log"
)

type BackupWorker struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
	opts    MetricsBackupOpts
	service *metricservice.Service
	uc      *metricusecase.MetricUseCase
}

func (w *BackupWorker) Save(ctx context.Context) error {
	metrics, err := w.findMetrics(ctx)
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return nil
	}
	if err = w.SaveToFile(metrics); err != nil {
		return err
	}

	return nil
}

func (w *BackupWorker) SaveToFile(metrics []entity.Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(w.opts.Filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}(file)
	writer := bufio.NewWriter(file)
	if _, err = writer.Write(data); err != nil {
		return err
	}
	if err = writer.WriteByte('\n'); err != nil {
		return err
	}

	return writer.Flush()
}

func (w *BackupWorker) Restore(ctx context.Context) error {
	if !w.scanner.Scan() {
		return w.scanner.Err()
	}
	data := w.scanner.Bytes()
	var metrics []entity.Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return err
	}
	for _, metric := range metrics {
		mv, err := metric.MetricValue()
		if err != nil {
			log.Error().Err(err).Send()
		}
		if _, err = w.uc.Update(ctx, metricusecase.MetricUpdateDTO{
			Type:  metric.MType,
			Name:  metric.ID,
			Value: mv,
		}); err != nil {
			log.Error().Err(err).Send()
		}
	}

	return nil
}

func (w *BackupWorker) Close() error {
	return w.file.Close()
}

type MetricsBackupOpts struct {
	Restore  bool
	Interval int
	Filepath string
}

func NewBackupWorker(
	service *metricservice.Service,
	uc *metricusecase.MetricUseCase,
	opts MetricsBackupOpts,
) (*BackupWorker, error) {
	file, err := os.OpenFile(opts.Filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &BackupWorker{
		service: service,
		uc:      uc,
		file:    file,
		scanner: bufio.NewScanner(file),
		writer:  bufio.NewWriter(file),
		opts:    opts,
	}, nil
}

func (w *BackupWorker) Run(ctx context.Context) error {
	if w.opts.Restore {
		if err := w.Restore(ctx); err != nil {
			log.Error().Err(err).Msg("failed to restore storage from file")
		}
		log.Info().Msg("successfully restored storage from file")
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error, 2)
	timer := time.NewTimer(time.Duration(w.opts.Interval) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			if err := w.Save(ctx); err != nil {
				log.Error().Err(err).Msg("failed to save storage to file")
			}

			return ctx.Err()
		case err := <-errCh:
			if err != nil {
				log.Error().Err(err).Msg("backup worker error")
			}
		case <-timer.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer timer.Reset(time.Duration(w.opts.Interval) * time.Second)

				errCh <- w.Save(ctx)
				log.Info().Msg("backup worked success")
			}()
		}
	}
}

func (w *BackupWorker) findMetrics(ctx context.Context) ([]entity.Metrics, error) {
	metricsDTO, err := w.service.FindAll(ctx, metricservice.FindAllDTO{})
	if err != nil {
		return nil, err
	}
	var metrics []entity.Metrics
	for _, dto := range metricsDTO {
		vl, aErr := anyToString(dto.Value)
		if aErr != nil {
			log.Error().Msgf("uncorrect convert to string: %v", dto.Value)
			vl = ""
		}
		m, bErr := entity.BuildMetrics(entity.MetricDTO{
			Type:  dto.Type,
			Name:  dto.Name,
			Value: vl,
		})
		if bErr != nil {
			continue
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func anyToString(value any) (string, error) {
	switch v := value.(type) {
	case int64:
		return strconv.Itoa(int(v)), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}
