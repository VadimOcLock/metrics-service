package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/VadimOcLock/metrics-service/internal/errorz"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
)

type UpdateMetricsHandler struct {
	ctx            context.Context
	MetricsUseCase metricusecase.UseCase
}

func NewUpdateMetricsHandler(
	ctx context.Context,
	uc metricusecase.UseCase,
) UpdateMetricsHandler {

	return UpdateMetricsHandler{
		ctx:            ctx,
		MetricsUseCase: uc,
	}
}

var _ http.Handler = (*UpdateMetricsHandler)(nil)

func (h UpdateMetricsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, errorz.ErrMsgOnlyPOSTMethodAccept, http.StatusMethodNotAllowed)

		return
	}
	parts := strings.Split(strings.TrimPrefix(req.URL.Path, "/update/"), "/")
	if len(parts) != 3 {
		http.Error(res, "invalid count input params", http.StatusBadRequest)

		return
	}
	for _, part := range parts {
		if part == "" {
			http.Error(res, errorz.ErrMsgEmptyMetricParam, http.StatusNotFound)

			return
		}
	}
	dto := entity.MetricDTO{
		Type:  parts[0],
		Name:  parts[1],
		Value: parts[2],
	}
	bodyObj, err := h.MetricsUseCase.UpdateMetric(h.ctx, dto)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)

		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(bodyObj)
	if err != nil {
		log.Printf("marshalling response body err: %s", err)

		return
	}
	_, err = res.Write(respBody)
	if err != nil {
		log.Printf("response body write err: %s", err)
	}
}
