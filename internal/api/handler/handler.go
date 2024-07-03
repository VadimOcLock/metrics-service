package handler

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/go-chi/chi/v5"
)

func New(ctx context.Context) http.Handler {
	store := somestore.New()
	metricService := metricservice.New(&store)
	metricUseCase := metricusecase.New(metricService)
	mh := NewMetricsHandler(ctx, metricUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", mh.GetAllMetrics)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", mh.UpdateMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", mh.GetMetricValue)
	})

	return r
}
