package metrichandler_test

import (
	"context"
	"errors"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/metrichandler"

	"github.com/VadimOcLock/metrics-service/internal/api/handlers/metrichandler/mocks"

	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

type updateMetricHandlerInput struct {
	method      string
	query       string
	metricType  string
	metricName  string
	metricValue string
}

type updateMetricHandlerWant struct {
	statusCode int
	response   string
}

type updateMetricHandlerTestCase struct {
	name  string
	input updateMetricHandlerInput
	want  updateMetricHandlerWant
}

func updateMetricHandlerTestCases() []updateMetricHandlerTestCase {
	tests := []updateMetricHandlerTestCase{
		{
			name: "non POST method",
			input: updateMetricHandlerInput{
				method: http.MethodGet,
				query:  "/update/",
			},
			want: updateMetricHandlerWant{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "missing parameters",
			input: updateMetricHandlerInput{
				method: http.MethodPost,
				query:  "/update/",
			},
			want: updateMetricHandlerWant{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "undefined metric type",
			input: updateMetricHandlerInput{
				method:      http.MethodPost,
				query:       "/update/",
				metricType:  "undefined_type",
				metricName:  "metric1",
				metricValue: "123.45",
			},
			want: updateMetricHandlerWant{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "successful gauge metric update",
			input: updateMetricHandlerInput{
				method:      http.MethodPost,
				query:       "/update/",
				metricType:  "gauge",
				metricName:  "metric1",
				metricValue: "123.45",
			},
			want: updateMetricHandlerWant{
				statusCode: http.StatusOK,
				response:   `{"message":"Metric updated successfully"}`,
			},
		},
	}

	return tests
}

func TestMetricsHandler_UpdateMetric(t *testing.T) {
	metricUseCase := mocks.NewMetricUseCase(t)
	h := metrichandler.NewMetricHandler(metricUseCase)

	tests := updateMetricHandlerTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricUseCase.ExpectedCalls = nil
			metricUseCase.Calls = nil

			r := httptest.NewRequest(tt.input.method, tt.input.query, nil)
			ctx := chi.NewRouteContext()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()

			switch tt.name {
			case "undefined metric type":
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					MType: "undefined_type",
					ID:    "metric1",
				}).Return(entity.Metrics{}, errorz.ErrUndefinedMetricType)
			case "undefined metric name":
				val := 123.45
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					MType: "gauge",
					ID:    "undefined_name",
					Value: &val,
				}).Return(entity.Metrics{}, errorz.ErrUndefinedMetricName)
			case "successful gauge metric update":
				val := 123.45
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					MType: "gauge",
					ID:    "metric1",
					Value: &val,
				}).Return(entity.Metrics{
					ID:    "metric1",
					MType: "gauge",
					Value: &val,
				}, nil)
			case "successful gauge metric update":
				val := 123.45
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					MType: "gauge",
					ID:    "metric1",
					Value: &val,
				}).Return(entity.Metrics{
					ID:    "metric1",
					MType: "gauge",
					Value: &val,
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

type getMetricHandlerInput struct {
	method     string
	query      string
	metricType string
	metricName string
}

type getMetricHandlerWant struct {
	statusCode int
	response   string
}

type getMetricHandlerTestCase struct {
	name  string
	input getMetricHandlerInput
	want  getMetricHandlerWant
}

func getMetricHandlerTestCases() []getMetricHandlerTestCase {
	tests := []getMetricHandlerTestCase{
		{
			name: "non GET method",
			input: getMetricHandlerInput{
				method: http.MethodPost,
				query:  "/value/",
			},
			want: getMetricHandlerWant{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "missing parameters",
			input: getMetricHandlerInput{
				method: http.MethodGet,
				query:  "/value/",
			},
			want: getMetricHandlerWant{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "undefined metric type",
			input: getMetricHandlerInput{
				method:     http.MethodGet,
				query:      "/value/undefined_type/metric1",
				metricType: "undefined_type",
				metricName: "metric1",
			},
			want: getMetricHandlerWant{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "undefined metric name",
			input: getMetricHandlerInput{
				method:     http.MethodGet,
				query:      "/value/gauge/undefined_name",
				metricType: "gauge",
				metricName: "undefined_name",
			},
			want: getMetricHandlerWant{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "successful metric retrieval",
			input: getMetricHandlerInput{
				method:     http.MethodGet,
				query:      "/value/gauge/metric1",
				metricType: "gauge",
				metricName: "metric1",
			},
			want: getMetricHandlerWant{
				statusCode: http.StatusOK,
				response:   "123.45",
			},
		},
		{
			name: "internal server error",
			input: getMetricHandlerInput{
				method:     http.MethodGet,
				query:      "/value/gauge/metric1",
				metricType: "gauge",
				metricName: "metric1",
			},
			want: getMetricHandlerWant{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	return tests
}

func TestMetricsHandler_GetMetricValue(t *testing.T) {
	metricUseCase := mocks.NewMetricUseCase(t)
	h := metrichandler.NewMetricHandler(metricUseCase)

	tests := getMetricHandlerTestCases()

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
				metricUseCase.On("Find", r.Context(), metricusecase.MetricFindDTO{
					MetricType: "undefined_type",
					MetricName: "metric1",
				}).Return(metricusecase.MetricFindResp{}, errorz.ErrUndefinedMetricType)
			case "undefined metric name":
				metricUseCase.On("Find", r.Context(), metricusecase.MetricFindDTO{
					MetricType: "gauge",
					MetricName: "undefined_name",
				}).Return(metricusecase.MetricFindResp{}, errorz.ErrUndefinedMetricName)
			case "internal server error":
				metricUseCase.On("Find", r.Context(), metricusecase.MetricFindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.MetricFindResp{}, errors.New("some error"))
			case "successful metric retrieval":
				metricUseCase.On("Find", r.Context(), metricusecase.MetricFindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.MetricFindResp{
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

type getAllMetricsHandlerInput struct {
	method string
	query  string
}

type getAllMetricsHandlerWant struct {
	statusCode int
	response   string
}

type getAllMetricsHandlerTestCase struct {
	name  string
	input getAllMetricsHandlerInput
	want  getAllMetricsHandlerWant
}

func getAllMetricsHandlerTestCases() []getAllMetricsHandlerTestCase {
	tests := []getAllMetricsHandlerTestCase{
		{
			name: "non GET method",
			input: getAllMetricsHandlerInput{
				method: http.MethodPost,
				query:  "/",
			},
			want: getAllMetricsHandlerWant{
				statusCode: http.StatusMethodNotAllowed,
			},
		},

		{
			name: "successful get all metrics",
			input: getAllMetricsHandlerInput{
				method: http.MethodGet,
				query:  "/",
			},
			want: getAllMetricsHandlerWant{
				statusCode: http.StatusOK,
				response:   "<html>success</html>",
			},
		},
		{
			name: "internal server error",
			input: getAllMetricsHandlerInput{
				method: http.MethodGet,
				query:  "/",
			},
			want: getAllMetricsHandlerWant{
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	return tests
}

func TestMetricsHandler_GetAllMetrics(t *testing.T) {
	metricUseCase := mocks.NewMetricUseCase(t)
	h := metrichandler.NewMetricHandler(metricUseCase)

	tests := getAllMetricsHandlerTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricUseCase.ExpectedCalls = nil
			metricUseCase.Calls = nil

			r := httptest.NewRequest(tt.input.method, tt.input.query, nil)
			w := httptest.NewRecorder()

			switch tt.name {
			case "successful get all metrics":
				metricUseCase.On("FindAll", r.Context(), metricusecase.MetricFindAllDTO{}).
					Return(metricusecase.MetricFindAllResp{
						HTML: "<html>success</html>",
					}, nil)
			case "internal server error":
				metricUseCase.On("FindAll", r.Context(), metricusecase.MetricFindAllDTO{}).
					Return(metricusecase.MetricFindAllResp{}, errors.New("some error"))
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
