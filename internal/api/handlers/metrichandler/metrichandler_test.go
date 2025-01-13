package metrichandler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/stretchr/testify/mock"

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
				query:  "/update/gauge/metric1/value",
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
				query:       "/update/undefined_type/metric1/123.45",
				metricType:  "undefined_type",
				metricName:  "metric1",
				metricValue: "123.45",
			},
			want: updateMetricHandlerWant{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "successful metric update",
			input: updateMetricHandlerInput{
				method:      http.MethodPost,
				query:       "/update/gauge/metric1/123.45",
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
			ctx.URLParams.Add("type", tt.input.metricType)
			ctx.URLParams.Add("name", tt.input.metricName)
			ctx.URLParams.Add("value", tt.input.metricValue)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()

			switch tt.name {
			case "undefined metric type":
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					Type:  "undefined_type",
					Name:  "metric1",
					Value: "123.45",
				}).Return(metricusecase.MetricUpdateResp{}, errorz.ErrUndefinedMetricType)
			case "undefined metric name":
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					Type:  "gauge",
					Name:  "undefined_name",
					Value: "123.45",
				}).Return(metricusecase.MetricUpdateResp{}, errorz.ErrUndefinedMetricName)
			case "successful metric update":
				metricUseCase.On("Update", r.Context(), metricusecase.MetricUpdateDTO{
					Type:  "gauge",
					Name:  "metric1",
					Value: "123.45",
				}).Return(metricusecase.MetricUpdateResp{
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

func TestMetricHandler_UpdateMetricJSON(t *testing.T) {
	type input struct {
		method string
		body   string
	}
	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name      string
		input     input
		want      want
		mockSetup func(useCase *mocks.MetricUseCase)
	}{
		{
			name: "non-POST method",
			input: input{
				method: http.MethodGet,
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				response:   errorz.ErrMsgOnlyPOSTMethodAccept + "\n",
			},
			mockSetup: nil,
		},
		{
			name: "decode error",
			input: input{
				method: http.MethodPost,
				body:   `invalid-json`,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   errorz.ErrInvalidRequestBody + "\n",
			},
			mockSetup: nil,
		},
		{
			name: "validation error",
			input: input{
				method: http.MethodPost,
				body:   `{"id":"test_metric","type":"unknown_metric_type"}`,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "unknown metric type: " + "unknown_metric_type\n",
			},
			mockSetup: func(useCase *mocks.MetricUseCase) {},
		},
		{
			name: "update error",
			input: input{
				method: http.MethodPost,
				body:   `{"id":"test_metric","type":"gauge","value":123.45}`,
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "update error\n",
			},
			mockSetup: func(useCase *mocks.MetricUseCase) {
				useCase.On("Update", mock.Anything, metricusecase.MetricUpdateDTO{
					Type:  "gauge",
					Name:  "test_metric",
					Value: "123.45",
				}).Return(metricusecase.MetricUpdateResp{}, errors.New("update error"))
			},
		},
		{
			name: "successful update",
			input: input{
				method: http.MethodPost,
				body:   `{"id":"test_metric","type":"gauge","value":123.45}`,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   `{"message":"Metric updated successfully"}`,
			},
			mockSetup: func(useCase *mocks.MetricUseCase) {
				useCase.On("Update", mock.Anything, metricusecase.MetricUpdateDTO{
					Type:  "gauge",
					Name:  "test_metric",
					Value: "123.45",
				}).Return(metricusecase.MetricUpdateResp{
					Message: "Metric updated successfully",
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := mocks.NewMetricUseCase(t)
			if tt.mockSetup != nil {
				tt.mockSetup(useCase)
			}

			handler := metrichandler.NewMetricHandler(useCase)

			req := httptest.NewRequest(tt.input.method, "/", bytes.NewBufferString(tt.input.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateMetricJSON(w, req)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestMetricHandler_Ping(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		mockSetup func(pool *mocks.Pool)
		wantCode  int
		wantBody  string
	}{
		{
			name:     "non-GET method",
			method:   http.MethodPost,
			wantCode: http.StatusMethodNotAllowed,
			wantBody: http.StatusText(http.StatusMethodNotAllowed) + "\n",
		},
		{
			name:     "pool is nil",
			method:   http.MethodGet,
			wantCode: http.StatusInternalServerError,
			wantBody: "database unavailable now\n",
		},
		{
			name:   "ping database error",
			method: http.MethodGet,
			mockSetup: func(pool *mocks.Pool) {
				pool.On("Ping", context.Background()).Return(errors.New("ping error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "database unavailable now\n",
		},
		{
			name:   "successful ping",
			method: http.MethodGet,
			mockSetup: func(pool *mocks.Pool) {
				pool.On("Ping", context.Background()).Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pool *mocks.Pool
			if tt.mockSetup != nil {
				pool = mocks.NewPool(t)
				tt.mockSetup(pool)
			}

			handler := metrichandler.MetricHandler{}
			pingHandler := handler.Ping(pool)

			req := httptest.NewRequest(tt.method, "/ping", nil)
			w := httptest.NewRecorder()

			pingHandler(w, req)

			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.wantCode, res.StatusCode)
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestMetricHandler_UpdateMetricBatch(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		mockBehavior func(useCase *mocks.MetricUseCase)
		wantCode     int
		wantBody     string
	}{
		{
			name:         "method not allowed",
			method:       http.MethodGet,
			body:         "",
			mockBehavior: func(useCase *mocks.MetricUseCase) {},
			wantCode:     http.StatusMethodNotAllowed,
			wantBody:     "Method Not Allowed\n",
		},
		{
			name:         "invalid JSON body",
			method:       http.MethodPost,
			body:         "{invalid-json",
			mockBehavior: func(useCase *mocks.MetricUseCase) {},
			wantCode:     http.StatusBadRequest,
			wantBody:     "invalid character 'i' looking for beginning of object key string\n",
		},
		{
			name:   "use case returns error",
			method: http.MethodPost,
			body: `[
				{"id":"metric1","type":"gauge","value":123.45},
				{"id":"metric2","type":"counter","delta":10}
			]`,
			mockBehavior: func(useCase *mocks.MetricUseCase) {
				vl := 123.45
				delta := int64(10)
				useCase.On("UpdateBatch", mock.Anything, metricusecase.MetricsUpdateBatchDTO{
					Data: &[]entity.Metrics{
						{ID: "metric1", MType: "gauge", Value: &vl},
						{ID: "metric2", MType: "counter", Delta: &delta},
					},
				}).Return(errors.New("update batch error"))
			},
			wantCode: http.StatusBadRequest,
			wantBody: "update batch error\n",
		},
		{
			name:   "successful update",
			method: http.MethodPost,
			body: `[
				{"id":"metric1","type":"gauge","value":123.45},
				{"id":"metric2","type":"counter","delta":10}
			]`,
			mockBehavior: func(useCase *mocks.MetricUseCase) {
				vl := 123.45
				delta := int64(10)
				useCase.On("UpdateBatch", mock.Anything, metricusecase.MetricsUpdateBatchDTO{
					Data: &[]entity.Metrics{
						{ID: "metric1", MType: "gauge", Value: &vl},
						{ID: "metric2", MType: "counter", Delta: &delta},
					},
				}).Return(nil)
			},
			wantCode: http.StatusOK,
			wantBody: "success update metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := mocks.NewMetricUseCase(t)
			tt.mockBehavior(useCase)

			h := metrichandler.NewMetricHandler(useCase)

			req := httptest.NewRequest(tt.method, "/update-batch", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.UpdateMetricBatch(w, req)

			res := w.Result()
			assert.Equal(t, tt.wantCode, res.StatusCode)

			body, _ := io.ReadAll(res.Body)
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestMetricHandler_GetMetricValueJSON(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		mockBehavior func(useCase *mocks.MetricUseCase)
		wantCode     int
		wantBody     string
	}{
		{
			name:   "method not allowed",
			method: http.MethodGet,
			body:   "",
			mockBehavior: func(useCase *mocks.MetricUseCase) {
			},
			wantCode: http.StatusMethodNotAllowed,
			wantBody: "Method Not Allowed\n",
		},
		{
			name:   "invalid JSON body",
			method: http.MethodPost,
			body:   "{invalid-json",
			mockBehavior: func(useCase *mocks.MetricUseCase) {
			},
			wantCode: http.StatusBadRequest,
			wantBody: "invalid character 'i' looking for beginning of object key string\n",
		},
		{
			name:   "metric not found",
			method: http.MethodPost,
			body:   `{"id":"metric1","type":"gauge"}`,
			mockBehavior: func(useCase *mocks.MetricUseCase) {
				useCase.On("Find", mock.Anything, metricusecase.MetricFindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.MetricFindResp{}, errorz.ErrMetricNotFound)
			},
			wantCode: http.StatusNotFound,
			wantBody: errorz.ErrMetricNotFound.Error() + "\n",
		},
		{
			name:   "internal server error",
			method: http.MethodPost,
			body:   `{"id":"metric1","type":"gauge"}`,
			mockBehavior: func(useCase *mocks.MetricUseCase) {
				useCase.On("Find", mock.Anything, metricusecase.MetricFindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.MetricFindResp{}, errors.New("internal server error"))
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "Internal Server Error\n",
		},
		{
			name:   "successful find",
			method: http.MethodPost,
			body:   `{"id":"metric1","type":"gauge"}`,
			mockBehavior: func(useCase *mocks.MetricUseCase) {
				expVl := 123.45
				expData := entity.Metrics{
					ID:    "metric1",
					MType: "gauge",
					Value: &expVl,
				}
				useCase.On("Find", mock.Anything, metricusecase.MetricFindDTO{
					MetricType: "gauge",
					MetricName: "metric1",
				}).Return(metricusecase.MetricFindResp{
					Data: &expData,
				}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":"metric1","type":"gauge","value":123.45}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := mocks.NewMetricUseCase(t)
			tt.mockBehavior(useCase)
			h := metrichandler.NewMetricHandler(useCase)

			req := httptest.NewRequest(tt.method, "/get-metric-value", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.GetMetricValueJSON(w, req)

			res := w.Result()
			assert.Equal(t, tt.wantCode, res.StatusCode)

			body, _ := io.ReadAll(res.Body)
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestGetMetricsValidateErr(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		want     bool
	}{
		{
			name:     "some expected error",
			inputErr: errorz.ErrUndefinedMetricType,
			want:     true,
		},
		{
			name:     "unrelated error",
			inputErr: errors.New("some unrelated error"),
			want:     false,
		},
		{
			name:     "nil error",
			inputErr: nil,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := metrichandler.GetMetricsValidateErr(tt.inputErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
