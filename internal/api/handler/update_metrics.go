package handler

import (
	"context"
	"encoding/json"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metric_usecase"
	"log"
	"net/http"
	"strings"
)

type UpdateMetricsHandler struct {
	ctx            context.Context
	MetricsUseCase metric_usecase.UseCase
}

func NewUpdateMetricsHandler(
	ctx context.Context,
	uc metric_usecase.UseCase,
) UpdateMetricsHandler {
	return UpdateMetricsHandler{
		ctx:            ctx,
		MetricsUseCase: uc,
	}
}

var _ http.Handler = (*UpdateMetricsHandler)(nil)

func (h UpdateMetricsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "only POST method accept", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(res, "invalid count input params", http.StatusBadRequest)
		return
	}
	dto := entity.MetricDTO{
		Type:  parts[2],
		Name:  parts[3],
		Value: parts[4],
	}
	bodyObj, err := h.MetricsUseCase.UpdateMetric(h.ctx, dto)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
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
