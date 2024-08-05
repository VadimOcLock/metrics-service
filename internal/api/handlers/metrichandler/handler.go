package metrichandler

import (
	"net/http"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/go-chi/chi/v5"
)

func New(
	store metricservice.Store,
	fileUpdater chan bool,
) http.Handler {
	metricService := metricservice.New(store)
	metricUseCase := metricusecase.New(&metricService)
	mh := NewMetricHandler(&metricUseCase, fileUpdater)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.GZIP)
	r.Use(chimiddleware.Recoverer)

	r.Get("/", mh.GetAllMetrics)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", mh.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", mh.UpdateMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", mh.GetMetricValueJSON)
		r.Get("/{type}/{name}", mh.GetMetricValue)
	})

	return r
}
