package handler

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/service/metric_service"
	"github.com/VadimOcLock/metrics-service/internal/store/some_store"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metric_usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func New(ctx context.Context) http.Handler {
	// Использовал chi так как не нашел способа распарсить пустые сегменты без регулярок
	r := chi.NewRouter()

	store := some_store.New()
	metricService := metric_service.New(&store)
	metricUseCase := metric_usecase.New(metricService)
	updateMetricsHandler := NewUpdateMetricsHandler(ctx, metricUseCase)

	r.Post("/update/{type:[^/]*}/{name:[^/]*}/{value:[^/]*}", updateMetricsHandler.ServeHTTP)

	return Chain(r, LoggerMiddleware)
}
