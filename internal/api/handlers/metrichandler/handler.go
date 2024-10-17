package metrichandler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func New(mh MetricHandler, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.GZipMiddleware)

	r.Get("/", mh.GetAllMetrics)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if pool == nil {
			http.Error(w, "database unavailable now", http.StatusInternalServerError)

			return
		}
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable now", http.StatusInternalServerError)

			return
		}
		w.WriteHeader(http.StatusOK)
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", mh.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", mh.UpdateMetric)
	})
	r.Post("/updates/", mh.UpdateMetricBatch)
	r.Route("/value", func(r chi.Router) {
		r.Post("/", mh.GetMetricValueJSON)
		r.Get("/{type}/{name}", mh.GetMetricValue)
	})

	return r
}
