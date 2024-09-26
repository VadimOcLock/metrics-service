package metrichandler

import (
	"net/http"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func New(mh MetricHandler) http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.GZipMiddleware)

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
