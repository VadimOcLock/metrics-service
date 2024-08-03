package metricusecase

import "github.com/VadimOcLock/metrics-service/internal/entity"

//type MetricUpdateResp struct {
//	Message string `json:"message"`
//}

//type MetricUpdateResp entity.Metrics

//type MetricUpdateDTO entity.MetricDTO

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

//type MetricFindResp struct {
//	MetricValue string
//}

//type MetricFindResp entity.Metrics
