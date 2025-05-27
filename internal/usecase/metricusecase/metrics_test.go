package metricusecase_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMetricUseCase_Update(t *testing.T) {
	expVl := 23.5
	expVlStr := fmt.Sprintf("%.1f", expVl)
	expCnt := int64(100)
	tests := []struct {
		name         string
		dto          metricusecase.MetricUpdateDTO
		mockSetup    func(service *mocks.MetricService)
		expectedResp metricusecase.MetricUpdateResp
		expectedErr  error
	}{
		{
			name: "Корректное обновление Gauge метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.GaugeMetricType,
				Name:  "test-gauge",
				Value: expVlStr,
			},
			mockSetup: func(service *mocks.MetricService) {
				service.On("UpdateGauge", context.Background(),
					metricservice.UpdateGaugeDTO{Name: "test-gauge", Value: expVl}).
					Return(nil)
			},
			expectedResp: metricusecase.MetricUpdateResp{
				Message: "metric update success",
				Data: &entity.Metrics{
					ID:    "test-gauge",
					MType: enum.GaugeMetricType,
					Value: &expVl,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Неверное значение для Gauge метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.GaugeMetricType,
				Name:  "test-gauge",
				Value: "invalid-value",
			},
			mockSetup:    func(service *mocks.MetricService) {},
			expectedResp: metricusecase.MetricUpdateResp{},
			expectedErr:  errorz.ErrInvalidMetricValue,
		},
		{
			name: "Пустое имя для Gauge метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.GaugeMetricType,
				Name:  "",
				Value: expVlStr,
			},
			mockSetup:    func(service *mocks.MetricService) {},
			expectedResp: metricusecase.MetricUpdateResp{},
			expectedErr:  errorz.ErrInvalidMetricName,
		},
		{
			name: "Ошибка при обновлении Gauge метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.GaugeMetricType,
				Name:  "test-gauge",
				Value: expVlStr,
			},
			mockSetup: func(service *mocks.MetricService) {
				service.On(
					"UpdateGauge",
					context.Background(),
					metricservice.UpdateGaugeDTO{Name: "test-gauge", Value: expVl}).
					Return(errors.New("update failed"))
			},
			expectedResp: metricusecase.MetricUpdateResp{},
			expectedErr:  errorz.ErrUpdateMetricFailed,
		},
		{
			name: "Корректное обновление Counter метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.CounterMetricType,
				Name:  "test-counter",
				Value: strconv.FormatInt(expCnt, 10),
			},
			mockSetup: func(service *mocks.MetricService) {
				service.On(
					"UpdateCounter",
					context.Background(),
					metricservice.UpdateCounterDTO{Name: "test-counter", Value: expCnt}).
					Return(nil)
			},
			expectedResp: metricusecase.MetricUpdateResp{
				Message: "metric update success",
				Data: &entity.Metrics{
					ID:    "test-counter",
					MType: enum.CounterMetricType,
					Delta: &expCnt,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Неверное значение для Counter метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.CounterMetricType,
				Name:  "test-counter",
				Value: "invalid-value",
			},
			mockSetup: func(service *mocks.MetricService) {
			},
			expectedResp: metricusecase.MetricUpdateResp{},
			expectedErr:  errorz.ErrInvalidMetricValue,
		},
		{
			name: "Пустое имя для Counter метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  enum.CounterMetricType,
				Name:  "",
				Value: strconv.FormatInt(expCnt, 10),
			},
			mockSetup: func(service *mocks.MetricService) {
			},
			expectedResp: metricusecase.MetricUpdateResp{},
			expectedErr:  errorz.ErrInvalidMetricName,
		},
		{
			name: "Неопределённый тип метрики",
			dto: metricusecase.MetricUpdateDTO{
				Type:  "UndefinedType", // Неопределённый тип
				Name:  "test-metric",
				Value: strconv.FormatInt(expCnt, 10),
			},
			mockSetup: func(service *mocks.MetricService) {
			},
			expectedResp: metricusecase.MetricUpdateResp{},
			expectedErr:  errorz.ErrUndefinedMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewMetricService(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}
			uc := metricusecase.New(mockService)
			resp, err := uc.Update(context.Background(), tt.dto)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}

func TestMetricUseCase_FindAll(t *testing.T) {
	expVl := 12.34
	cntVl := int64(100)
	tests := []struct {
		name          string
		mockSetup     func(mockMetricService *mocks.MetricService, mockHTMLBuilder *mocks.HTMLBuilder)
		expectedResp  metricusecase.MetricFindAllResp
		expectedError error
	}{
		{
			name: "Successful FindAllWithHTML",
			mockSetup: func(mockMetricService *mocks.MetricService, mockHTMLBuilder *mocks.HTMLBuilder) {
				mockMetricService.On(
					"FindAll",
					context.Background(),
					metricservice.FindAllDTO{}).
					Return([]entity.Metrics{
						{ID: "test-gauge", MType: enum.GaugeMetricType, Value: &expVl},
						{ID: "test-cnt", MType: enum.CounterMetricType, Delta: &cntVl},
					}, nil)

				mockHTMLBuilder.On("BuildHTML", []entity.Metrics{
					{ID: "test-gauge", MType: enum.GaugeMetricType, Value: &expVl},
					{ID: "test-cnt", MType: enum.CounterMetricType, Delta: &cntVl},
				}).Return("<html>...</html>", nil)
			},
			expectedResp: metricusecase.MetricFindAllResp{
				HTML: "<html>...</html>",
			},
			expectedError: nil,
		},
		{
			name: "Error in MetricService FindAllWithHTML",
			mockSetup: func(mockMetricService *mocks.MetricService, mockHTMLBuilder *mocks.HTMLBuilder) {
				mockMetricService.On(
					"FindAll",
					context.Background(),
					metricservice.FindAllDTO{}).
					Return(nil, errors.New("service error"))
			},
			expectedResp:  metricusecase.MetricFindAllResp{},
			expectedError: fmt.Errorf("metricusecase.FindAllWithHTML: %w", errors.New("service error")),
		},
		{
			name: "Error in HTMLBuilder BuildHTML",
			mockSetup: func(mockMetricService *mocks.MetricService, mockHTMLBuilder *mocks.HTMLBuilder) {
				mockMetricService.On(
					"FindAll",
					context.Background(),
					metricservice.FindAllDTO{}).
					Return([]entity.Metrics{
						{ID: "test-gauge", MType: enum.GaugeMetricType, Value: &expVl},
					}, nil)

				mockHTMLBuilder.On("BuildHTML", []entity.Metrics{
					{ID: "test-gauge", MType: enum.GaugeMetricType, Value: &expVl},
				}).Return("", errors.New("html error"))
			},
			expectedResp:  metricusecase.MetricFindAllResp{},
			expectedError: fmt.Errorf("metricusecase.FindAllWithHTML: %w", errors.New("html error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMetricService := mocks.NewMetricService(t)
			mockHTMLBuilder := mocks.NewHTMLBuilder(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockMetricService, mockHTMLBuilder)
			}

			uc := metricusecase.New(mockMetricService, metricusecase.WithHTMLBuilder(mockHTMLBuilder))
			resp, err := uc.FindAllWithHTML(context.Background(), metricusecase.MetricFindAllDTO{})

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}

func TestMetricUseCase_Find(t *testing.T) {
	expVl := 12.34
	expVlStr := fmt.Sprintf("%.2f", expVl)
	tests := []struct {
		name          string
		dto           metricusecase.MetricFindDTO
		mockSetup     func(mockMetricService *mocks.MetricService)
		expectedResp  metricusecase.MetricFindResp
		expectedError error
	}{
		{
			name: "Successful Find",
			dto: metricusecase.MetricFindDTO{
				MetricType: enum.GaugeMetricType,
				MetricName: "test-gauge",
			},
			mockSetup: func(mockMetricService *mocks.MetricService) {
				mockMetricService.On("Find", context.Background(), metricservice.FindDTO{
					MetricType: enum.GaugeMetricType,
					MetricName: "test-gauge",
				}).Return(entity.Metrics{
					ID:    "test-gauge",
					MType: enum.GaugeMetricType,
					Value: &expVl,
				}, nil)
			},
			expectedResp: metricusecase.MetricFindResp{
				MetricValue: expVlStr,
				Data: &entity.Metrics{
					ID:    "test-gauge",
					MType: enum.GaugeMetricType,
					Value: &expVl,
				},
			},
			expectedError: nil,
		},
		{
			name: "Error in MetricService Find",
			dto: metricusecase.MetricFindDTO{
				MetricType: enum.GaugeMetricType,
				MetricName: "test-gauge",
			},
			mockSetup: func(mockMetricService *mocks.MetricService) {
				mockMetricService.On("Find", context.Background(), metricservice.FindDTO{
					MetricType: enum.GaugeMetricType,
					MetricName: "test-gauge",
				}).Return(entity.Metrics{}, errors.New("service error"))
			},
			expectedResp:  metricusecase.MetricFindResp{},
			expectedError: fmt.Errorf("metricusecase.Find: %w", errors.New("service error")),
		},
		{
			name: "Error in MetricValue",
			dto: metricusecase.MetricFindDTO{
				MetricType: "unknown_metric_type",
				MetricName: "test-gauge",
			},
			mockSetup: func(mockMetricService *mocks.MetricService) {
				mockMetricService.On("Find", context.Background(), metricservice.FindDTO{
					MetricType: "unknown_metric_type",
					MetricName: "test-gauge",
				}).Return(entity.Metrics{
					ID:    "test-gauge",
					MType: "unknown_metric_type",
				}, nil)
			},
			expectedResp:  metricusecase.MetricFindResp{},
			expectedError: errors.New("metricusecase.Find: unknown metric type: unknown_metric_type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMetricService := mocks.NewMetricService(t)

			if tt.mockSetup != nil {
				tt.mockSetup(mockMetricService)
			}

			uc := metricusecase.New(mockMetricService, metricusecase.WithHTMLBuilder(mocks.NewHTMLBuilder(t)))
			resp, err := uc.Find(context.Background(), tt.dto)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
