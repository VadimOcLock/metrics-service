package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const errMsgInvalidMetricType = "invalid metric type"
const errMsgInvalidMetricValue = "invalid metric value"

type MockUseCase struct{}

func TestNewUpdateMetricsHandler(t *testing.T) {
	type input struct {
		ctx context.Context
		uc  metricusecase.UseCase
	}
	tests := []struct {
		name  string
		input input
		want  MetricsHandler
	}{
		{
			name: "Valid context and use case",
			input: input{
				ctx: context.Background(),
				uc:  &MockUseCase{},
			},
			want: MetricsHandler{
				ctx:            context.Background(),
				MetricsUseCase: &MockUseCase{},
			},
		},
		{
			name: "Nil context",
			input: input{
				ctx: nil,
				uc:  &MockUseCase{},
			},
			want: MetricsHandler{
				ctx:            nil,
				MetricsUseCase: &MockUseCase{},
			},
		},
		{
			name: "Nil use case",
			input: input{
				ctx: context.Background(),
				uc:  nil,
			},
			want: MetricsHandler{
				ctx:            context.Background(),
				MetricsUseCase: nil,
			},
		},
		{
			name: "Both context and use case nil",
			input: input{
				ctx: nil,
				uc:  nil,
			},
			want: MetricsHandler{
				ctx:            nil,
				MetricsUseCase: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricsHandler(tt.input.ctx, tt.input.uc)

			assert.Equal(t, tt.want.ctx, h.ctx)
			assert.Equal(t, tt.want.MetricsUseCase, h.MetricsUseCase)
		})
	}
}

func (m *MockUseCase) Update(_ context.Context, dto entity.MetricDTO) (metricusecase.UpdateResp, error) {
	if dto.Type == "invalid" {

		return metricusecase.UpdateResp{}, errors.New(errMsgInvalidMetricType)
	}
	if dto.Value == "invalid" {

		return metricusecase.UpdateResp{}, errors.New(errMsgInvalidMetricValue)
	}

	return metricusecase.UpdateResp{Message: "success"}, nil
}

func TestUpdateMetricsHandler_ServeHTTP(t *testing.T) {
	type (
		input struct {
			method string
			query  string
		}
		want struct {
			code        int
			response    string
			contentType string
		}
	)

	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "correct input params",
			input: input{
				method: http.MethodPost,
				query:  "/update/valid_type/valid_name/valid_value",
			},
			want: want{
				code:        http.StatusOK,
				response:    `{"message":"success"}`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "invalid HTTP method",
			input: input{
				method: http.MethodGet,
				query:  "/update/valid_type/valid_name/valid_value",
			},
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    errorz.ErrMsgOnlyPOSTMethodAccept,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "invalid input metric type",
			input: input{
				method: http.MethodPost,
				query:  "/update/invalid/valid_name/valid_value",
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    errMsgInvalidMetricType,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "invalid input metric value",
			input: input{
				method: http.MethodPost,
				query:  "/update/valid_type/valid_name/invalid",
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    errMsgInvalidMetricValue,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "empty input metric name",
			input: input{
				method: http.MethodPost,
				query:  "/update//valid_name/valid_value",
			},
			want: want{
				code:        http.StatusNotFound,
				response:    errorz.ErrMsgEmptyMetricParam,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.input.method, tt.input.query, nil)
			w := httptest.NewRecorder()

			handler := MetricsHandler{
				ctx:            context.Background(),
				MetricsUseCase: &MockUseCase{},
			}
			handler.UpdateMetric(w, req)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			defer func() {
				_ = res.Body.Close()
			}()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, strings.TrimSuffix(string(resBody), "\n"))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
