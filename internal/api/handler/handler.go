package handler

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/service/metric_service"
	"github.com/VadimOcLock/metrics-service/internal/store/some_store"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metric_usecase"
	"net/http"
)

func New(ctx context.Context) http.Handler {
	mux := http.NewServeMux()

	store := some_store.New()
	metricService := metric_service.New(&store)
	metricUseCase := metric_usecase.New(metricService)
	updateMetricsHandler := NewUpdateMetricsHandler(ctx, metricUseCase)

	mux.Handle("/update/", updateMetricsHandler)

	return Chain(mux, LoggerMiddleware)
}
