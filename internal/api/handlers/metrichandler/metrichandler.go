package metrichandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"

	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/go-chi/chi/v5"

	"github.com/VadimOcLock/metrics-service/internal/entity"
)

// MetricHandler обрабатывает HTTP-запросы для работы с метриками.
// Содержит use case для бизнес-логики работы с метриками.
type MetricHandler struct {
	MetricsUseCase MetricUseCase
}

// Pool представляет интерфейс для проверки соединения с базой данных.
type Pool interface {
	// Ping проверяет соединение с базой данных.
	Ping(ctx context.Context) error
}

var _ MetricUseCase = (*metricusecase.MetricUseCase)(nil)

// NewMetricHandler создает новый экземпляр MetricHandler.
// Принимает:
//   - uc: реализацию интерфейса MetricUseCase
//
// Возвращает:
//   - новый экземпляр MetricHandler
func NewMetricHandler(
	uc MetricUseCase,
) MetricHandler {
	return MetricHandler{
		MetricsUseCase: uc,
	}
}

// UpdateMetric обрабатывает POST запрос для обновления метрики через URL параметры.
// Формат URL: /update/{type}/{name}/{value}
// Поддерживаемые типы метрик: gauge, counter.
func (h *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAccept, http.StatusMethodNotAllowed)

		return
	}
	dto := entity.MetricDTO{
		Type:  chi.URLParam(req, "type"),
		Name:  chi.URLParam(req, "name"),
		Value: chi.URLParam(req, "value"),
	}
	if dto.Type == "" || dto.Name == "" || dto.Value == "" {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}
	bodyObj, err := h.MetricsUseCase.Update(req.Context(), metricusecase.MetricUpdateDTO{
		Type:  dto.Type,
		Name:  dto.Name,
		Value: dto.Value,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(bodyObj)
	if err != nil {
		log.Error().Msgf("marshalling response body err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	if _, err = res.Write(respBody); err != nil {
		log.Error().Msgf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// UpdateMetricJSON обрабатывает POST запрос для обновления метрики через JSON тело.
// Пример тела запроса:
//
//	{
//	    "id": "Alloc",
//	    "type": "gauge",
//	    "value": 208256
//	}
func (h *MetricHandler) UpdateMetricJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAccept, http.StatusMethodNotAllowed)

		return
	}

	var dto entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		log.Error().Msgf("decode err: %s", err)
		http.Error(res, errorz.ErrInvalidRequestBody, http.StatusBadRequest)

		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err)
		}
	}(req.Body)
	if err := dto.Valid(); err != nil {
		log.Error().Msgf("valid err: %s", err)
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	mv, err := dto.MetricValue()
	if err != nil {
		log.Error().Msgf("dto err: %s", err)
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}

	bodyObj, err := h.MetricsUseCase.Update(req.Context(), metricusecase.MetricUpdateDTO{
		Type:  dto.MType,
		Name:  dto.ID,
		Value: mv,
	})
	if err != nil {
		log.Error().Msgf("update err: %s", err)
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(bodyObj)
	if err != nil {
		log.Error().Msgf("marshalling response body err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	if _, err = res.Write(respBody); err != nil {
		log.Error().Msgf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// GetAllMetrics возвращает HTML страницу со всеми текущими метриками.
func (h *MetricHandler) GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	r, err := h.MetricsUseCase.FindAllWithHTML(req.Context(), metricusecase.MetricFindAllDTO{})
	if err != nil {
		log.Error().Msgf("find all metrics err: %s", err)
		http.Error(res, errorz.ErrMsgFindAllMetrics, http.StatusInternalServerError)

		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	if _, err = res.Write([]byte(r.HTML)); err != nil {
		log.Error().Msgf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// GetMetricValue возвращает значение метрики в текстовом формате.
// Формат URL: /value/{type}/{name}
func (h *MetricHandler) GetMetricValue(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	metricType := chi.URLParam(req, "type")
	metricName := chi.URLParam(req, "name")
	if metricType == "" || metricName == "" {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}
	find, err := h.MetricsUseCase.Find(req.Context(), metricusecase.MetricFindDTO{
		MetricType: metricType,
		MetricName: metricName,
	})
	if GetMetricsValidateErr(err) {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}
	if err != nil {
		log.Error().Msgf("find metric err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	if _, err = res.Write([]byte(find.MetricValue)); err != nil {
		log.Error().Msgf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// GetMetricValueJSON возвращает значение метрики в JSON формате.
// Пример тела запроса:
//
//	{
//	    "id": "Alloc",
//	    "type": "gauge"
//	}
func (h *MetricHandler) GetMetricValueJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Debug().Msg("http.StatusMethodNotAllowed")
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	var dto entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		log.Error().Err(err).Send()
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err)
		}
	}(req.Body)

	find, err := h.MetricsUseCase.Find(req.Context(), metricusecase.MetricFindDTO{
		MetricType: dto.MType,
		MetricName: dto.ID,
	})
	if GetMetricsValidateErr(err) {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}
	if err != nil {
		log.Error().Msgf("find metric err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(find.Data)
	if err != nil {
		log.Error().Msgf("marshalling response body err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	if _, err = res.Write(respBody); err != nil {
		log.Error().Msgf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// GetMetricsValidateErr проверяет, является ли ошибка ошибкой валидации метрики.
func GetMetricsValidateErr(err error) bool {
	return errors.Is(err, errorz.ErrUndefinedMetricType) ||
		errors.Is(err, errorz.ErrUndefinedMetricName) ||
		errors.Is(err, errorz.ErrGaugeTypeNilValue) ||
		errors.Is(err, errorz.ErrCounterTypeNilDelta) ||
		errors.Is(err, errorz.ErrMetricNotFound)
}

// UpdateMetricBatch обрабатывает пакетное обновление метрик.
// Пример тела запроса:
// [
//
//	{
//	    "id": "Alloc",
//	    "type": "gauge",
//	    "value": 123.45
//	},
//	{
//	    "id": "PollCount",
//	    "type": "counter",
//	    "delta": 1
//	}
//
// ]
func (h *MetricHandler) UpdateMetricBatch(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Debug().Msg("http.StatusMethodNotAllowed")
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	var metrics []entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&metrics); err != nil {
		log.Error().Err(err).Send()
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)
	if err := h.MetricsUseCase.UpdateBatch(req.Context(),
		metricusecase.MetricsUpdateBatchDTO{Data: &metrics}); err != nil {
		log.Error().Msgf("update metrics err: %s", err)
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte("success update metrics")); err != nil {
		log.Error().Msgf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// Ping возвращает обработчик для проверки соединения с базой данных.
// Доступно по GET /ping.
func (h *MetricHandler) Ping(pool Pool) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Debug().Msg("http.StatusMethodNotAllowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}
		if pool == nil || reflect.ValueOf(pool).IsNil() {
			http.Error(w, "database unavailable now", http.StatusInternalServerError)

			return
		}
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable now", http.StatusInternalServerError)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
