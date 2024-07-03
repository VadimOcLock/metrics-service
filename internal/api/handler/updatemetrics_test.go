package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase/mocks"
	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

const errMsgInvalidMetricType = "invalid metric type"
const errMsgInvalidMetricValue = "invalid metric value"

func TestMetricsHandler_UpdateMetric(t *testing.T) {
	metricUseCase := mocks.NewUseCase(t)
	h := NewMetricsHandler(metricUseCase)

	type (
		input struct {
			method      string
			query       string
			metricType  string
			metricName  string
			metricValue string
		}
		want struct {
			statusCode int
			response   string
		}
	)

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "non POST method",
			input: input{
				method: http.MethodGet,
				query:  "/update/gauge/metric1/value",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "missing parameters",
			input: input{
				method: http.MethodPost,
				query:  "/update/",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "undefined metric type",
			input: input{
				method:      http.MethodPost,
				query:       "/update/undefined_type/metric1/123.45",
				metricType:  "undefined_type",
				metricName:  "metric1",
				metricValue: "123.45",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "successful metric update",
			input: input{
				method:      http.MethodPost,
				query:       "/update/gauge/metric1/123.45",
				metricType:  "gauge",
				metricName:  "metric1",
				metricValue: "123.45",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   `{"message":"Metric updated successfully"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricUseCase.ExpectedCalls = nil
			metricUseCase.Calls = nil

			r := httptest.NewRequest(tt.input.method, tt.input.query, nil)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("type", tt.input.metricType)
			ctx.URLParams.Add("name", tt.input.metricName)
			ctx.URLParams.Add("value", tt.input.metricValue)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()

			switch tt.name {
			case "undefined metric type":
				metricUseCase.On("Update", r.Context(), entity.MetricDTO{
					Type:  "undefined_type",
					Name:  "metric1",
					Value: "123.45",
				}).Return(metricusecase.UpdateResp{}, errorz.ErrUndefinedMetricType)
			case "undefined metric name":
				metricUseCase.On("Update", r.Context(), entity.MetricDTO{
					Type:  "gauge",
					Name:  "undefined_name",
					Value: "123.45",
				}).Return(metricusecase.UpdateResp{}, errorz.ErrUndefinedMetricName)
			case "successful metric update":
				metricUseCase.On("Update", r.Context(), entity.MetricDTO{
					Type:  "gauge",
					Name:  "metric1",
					Value: "123.45",
				}).Return(metricusecase.UpdateResp{
					Message: "Metric updated successfully",
				}, nil)
			}

			h.UpdateMetric(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			if tt.want.response != "" {
				body, _ := io.ReadAll(res.Body)
				defer func() {
					_ = res.Body.Close()
				}()
				assert.Equal(t, tt.want.response, string(body))
			}
		})
	}
}

func TestMetricsHandler_GetMetricValue(t *testing.T) {
	metricUseCase := mocks.NewUseCase(t)
	h := NewMetricsHandler(metricUseCase)

	type (
		input struct {
			method     string
			query      string
			metricType string
			metricName string
		}
		want struct {
			statusCode int
			response   string
		}
	)

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "non GET method",
			input: input{
				method: http.MethodPost,
				query:  "/value/",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "missing parameters",
			input: input{
				method: http.MethodGet,
				query:  "/value/",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "undefined metric type",
			input: input{
				method:     http.MethodGet,
				query:      "/value/undefined_type/metric1",
				metricType: "undefined_type",
				metricName: "metric1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "undefined metric name",
			input: input{
				method:     http.MethodGet,
				query:      "/value/gauge/undefined_name",
				metricType: "gauge",
				metricName: "undefined_name",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "successful metric retrieval",
			input: input{
				method:     http.MethodGet,
				query:      "/value/gauge/metric1",
				metricType: "gauge",
				metricName: "metric1",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   "123.45",
			},
		},
		{
			name: "internal server error",
			input: input{
				method:     http.MethodGet,
				query:      "/value/gauge/metric1",
				metricType: "gauge",
				metricName: "metric1",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricUseCase.ExpectedCalls = nil
			metricUseCase.Calls = nil

			r := httptest.NewRequest(tt.input.method, tt.input.query, nil)
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("type", tt.input.metricType)
			ctx.URLParams.Add("name", tt.input.metricName)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()

			switch tt.name {
			case "undefined metric type":
				metricUseCase.On("Find", r.Context(), metricusecase.FindDTO{
					MetricType: "undefined_type",
					MetricName: "metric1",
				}).Return(metricusecase.FindResp{}, errorz.ErrUndefinedMetricType)
			case "undefined metric name":
				metricUseCase.On("Find", r.Context(), metricusecase.FindDTO{
					MetricType: "gauge",
					MetricName: "undefined_name",
				}).Return(metricusecase.FindResp{}, errorz.ErrUndefinedMetricName)
			case "internal server error":
				metricUseCase.On("Find", r.Context(), metricusecase.FindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.FindResp{}, errors.New("some error"))
			case "successful metric retrieval":
				metricUseCase.On("Find", r.Context(), metricusecase.FindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.FindResp{
					MetricValue: "123.45",
				}, nil)
			}

			h.GetMetricValue(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			if tt.want.response != "" {
				body, _ := io.ReadAll(res.Body)
				defer func() {
					_ = res.Body.Close()
				}()
				assert.Equal(t, tt.want.response, string(body))
			}
		})
	}
}

func TestMetricsHandler_GetAllMetrics(t *testing.T) {
	metricUseCase := mocks.NewUseCase(t)
	h := NewMetricsHandler(metricUseCase)

	type (
		input struct {
			method string
			query  string
		}
		want struct {
			statusCode int
			response   string
		}
	)

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "non GET method",
			input: input{
				method: http.MethodPost,
				query:  "/",
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},

		{
			name: "successful get all metrics",
			input: input{
				method: http.MethodGet,
				query:  "/",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   "<html>success</html>",
			},
		},
		{
			name: "internal server error",
			input: input{
				method: http.MethodGet,
				query:  "/",
			},
			want: want{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricUseCase.ExpectedCalls = nil
			metricUseCase.Calls = nil

			r := httptest.NewRequest(tt.input.method, tt.input.query, nil)
			w := httptest.NewRecorder()

			switch tt.name {
			case "successful get all metrics":
				metricUseCase.On("FindAll", r.Context(), metricusecase.FindAllDTO{}).
					Return(metricusecase.FindAllResp{
						HTML: "<html>success</html>",
					}, nil)
			case "internal server error":
				metricUseCase.On("FindAll", r.Context(), metricusecase.FindAllDTO{}).
					Return(metricusecase.FindAllResp{}, errors.New("some error"))
			}

			h.GetAllMetrics(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			if tt.want.response != "" {
				body, _ := io.ReadAll(res.Body)
				defer func() {
					_ = res.Body.Close()
				}()
				assert.Equal(t, tt.want.response, string(body))
			}
		})
	}
}
