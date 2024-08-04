package metricusecase

import "github.com/VadimOcLock/metrics-service/internal/entity"

type MetricUpdateDTO entity.Metrics

type MetricFindAllResp struct {
	HTML string
}

type MetricFindAllDTO struct {
}

type MetricFindDTO struct {
	MetricType string
	MetricName string
}
