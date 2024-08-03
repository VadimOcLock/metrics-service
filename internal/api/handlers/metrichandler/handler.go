package metrichandler

import (
	"net/http"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/somestore"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	store := somestore.New()
	metricService := metricservice.New(&store)
	metricUseCase := metricusecase.New(&metricService)
	mh := NewMetricHandler(&metricUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)

	//r.Get("/", mh.GetAllMetrics)
	//r.Route("/update", func(r chi.Router) {
	//	r.Post("/{type}/{name}/{value}", mh.UpdateMetric)
	//})
	//r.Route("/value", func(r chi.Router) {
	//	r.Get("/{type}/{name}", mh.GetMetricValue)
	//})

	r.Get("/", mh.GetAllMetrics)
	r.Post("/update/", mh.UpdateMetric)
	r.Get("/value/", mh.GetMetricValue)

	return r
}
