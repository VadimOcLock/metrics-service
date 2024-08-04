package metrichandler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
)

type MetricHandler struct {
	MetricsUseCase MetricUseCase
}

var _ MetricUseCase = (*metricusecase.MetricUseCase)(nil)

func NewMetricHandler(
	uc MetricUseCase,
) MetricHandler {
	return MetricHandler{
		MetricsUseCase: uc,
	}
}

func (h *MetricHandler) GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	r, err := h.MetricsUseCase.FindAll(req.Context(), metricusecase.MetricFindAllDTO{})
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

func (h *MetricHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAccept, http.StatusMethodNotAllowed)

		return
	}
	var dto entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err)
		}
	}(req.Body)
	if err := dto.Valid(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	bodyObj, err := h.MetricsUseCase.Update(req.Context(), metricusecase.MetricUpdateDTO(dto))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	respBody, err := json.MarshalIndent(bodyObj, "", "    ")
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

func (h *MetricHandler) GetMetricValue(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	var dto entity.Metrics
	if err := json.NewDecoder(req.Body).Decode(&dto); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err)
		}
	}(req.Body)
	if err := dto.Valid(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	metric, err := h.MetricsUseCase.Find(req.Context(), metricusecase.MetricFindDTO{
		MetricType: dto.MType,
		MetricName: dto.ID,
	})
	if errors.Is(err, errorz.ErrUndefinedMetricType) ||
		errors.Is(err, errorz.ErrUndefinedMetricName) {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}
	if err != nil {
		log.Error().Msgf("find metric err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	respBody, err := json.MarshalIndent(metric, "", "    ")
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
