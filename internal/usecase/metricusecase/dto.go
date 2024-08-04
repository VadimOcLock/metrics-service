package metricusecase

import "github.com/VadimOcLock/metrics-service/internal/entity"

type MetricUpdateResp struct {
	Message string          `json:"message"`
	Data    *entity.Metrics `json:"-"`
}

type MetricUpdateDTO entity.MetricDTO

type MetricFindAllResp struct {
	HTML string
}

type MetricFindAllDTO struct {
}

type MetricFindDTO struct {
	MetricType string
	MetricName string
}

type MetricFindResp struct {
	MetricValue string
	Data        *entity.Metrics `json:"-"`
}
