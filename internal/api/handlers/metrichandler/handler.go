package metrichandler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

type HandlerConfig struct {
	DbPool             *pgxpool.Pool
	SecretSignatureKey string
}

type Option interface {
	apply(*HandlerConfig)
}

type optionFunc func(*HandlerConfig)

func (o optionFunc) apply(c *HandlerConfig) {
	o(c)
}

// WithSecretSignatureKey возвращает опцию для установки секретного ключа подписи в конфигурации HandlerConfig.
// Используйте эту функцию, чтобы задать секретный ключ при создании обработчика.
//
// Пример:
//
//	handler := New(myMetricHandler, WithSecretSignatureKey("my_secret_key"))
func WithSecretSignatureKey(key string) Option {
	return optionFunc(func(cfg *HandlerConfig) {
		cfg.SecretSignatureKey = key
	})
}

// WithDbPool возвращает опцию для установки пула подключений к базе данных в конфигурации HandlerConfig.
// Используйте эту функцию, чтобы передать пул подключений при создании обработчика.
//
// Пример:
//
//	handler := New(myMetricHandler, WithDbPool(myDbPool))
func WithDbPool(dbPool *pgxpool.Pool) Option {
	return optionFunc(func(cfg *HandlerConfig) {
		cfg.DbPool = dbPool
	})
}

func New(mh MetricHandler, opts ...Option) http.Handler {
	cfg := &HandlerConfig{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.GZipMiddleware)
	r.Use(middleware.RequestSignatureVerificationMiddleware(cfg.SecretSignatureKey))
	r.Use(middleware.ResponseSigningMiddleware(cfg.SecretSignatureKey))

	r.Get("/", mh.GetAllMetrics)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if cfg.DbPool == nil {
			http.Error(w, "database unavailable now", http.StatusInternalServerError)

			return
		}
		if err := cfg.DbPool.Ping(r.Context()); err != nil {
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
