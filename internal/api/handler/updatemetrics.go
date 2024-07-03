package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/go-chi/chi/v5"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
)

type MetricsHandler struct {
	MetricsUseCase metricusecase.UseCase
}

func NewMetricsHandler(
	uc metricusecase.UseCase,
) MetricsHandler {

	return MetricsHandler{
		MetricsUseCase: uc,
	}
}

func (h MetricsHandler) UpdateMetric(res http.ResponseWriter, req *http.Request) {
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
	bodyObj, err := h.MetricsUseCase.Update(req.Context(), dto)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(bodyObj)
	if err != nil {
		log.Printf("marshalling response body err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	if _, err = res.Write(respBody); err != nil {
		log.Printf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h MetricsHandler) GetAllMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	r, err := h.MetricsUseCase.FindAll(req.Context(), metricusecase.FindAllDTO{})
	if err != nil {
		log.Printf("find all metrics err: %s", err)
		http.Error(res, errorz.ErrMsgFindAllMetrics, http.StatusInternalServerError)

		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	if _, err = res.Write([]byte(r.HTML)); err != nil {
		log.Printf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (h MetricsHandler) GetMetricValue(res http.ResponseWriter, req *http.Request) {
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
	find, err := h.MetricsUseCase.Find(req.Context(), metricusecase.FindDTO{
		MetricType: metricType,
		MetricName: metricName,
	})
	if errors.Is(err, errorz.ErrUndefinedMetricType) ||
		errors.Is(err, errorz.ErrUndefinedMetricName) {
		http.Error(res, err.Error(), http.StatusNotFound)

		return
	}
	if err != nil {
		log.Printf("find metric err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	if _, err = res.Write([]byte(find.MetricValue)); err != nil {
		log.Printf("response body write err: %s", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}
