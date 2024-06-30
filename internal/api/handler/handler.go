package handler

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metric_usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func New(ctx context.Context) http.Handler {
	// Использовал chi так как не нашел способа распарсить пустые сегменты без регулярок
	r := chi.NewRouter()

	store := somestore.New()
	metricService := metricservice.New(&store)
	metricUseCase := metric_usecase.New(metricService)
	updateMetricsHandler := NewUpdateMetricsHandler(ctx, metricUseCase)

	r.Post("/update/{type:[^/]*}/{name:[^/]*}/{value:[^/]*}", updateMetricsHandler.ServeHTTP)

	return Chain(r, LoggerMiddleware)
}
