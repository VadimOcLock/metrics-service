package metricservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/service/metricservice/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_UpdateGauge(t *testing.T) {
	errStore := errors.New("store error")
	tests := []struct {
		name        string
		dto         metricservice.UpdateGaugeDTO
		mockSetup   func(store *mocks.Store)
		expectedErr error
	}{
		{
			name: "Success update",
			dto: metricservice.UpdateGaugeDTO{
				Name:  "test-gauge",
				Value: 12.34,
			},
			mockSetup: func(store *mocks.Store) {
				store.On("UpsertGaugeMetric", mock.Anything, metricservice.UpsertGaugeMetricParams{
					Name:  "test-gauge",
					Value: 12.34,
				}).Return(true, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Store error",
			dto: metricservice.UpdateGaugeDTO{
				Name:  "test-gauge",
				Value: 12.34,
			},
			mockSetup: func(store *mocks.Store) {
				store.On("UpsertGaugeMetric", mock.Anything, metricservice.UpsertGaugeMetricParams{
					Name:  "test-gauge",
					Value: 12.34,
				}).Return(false, errStore)
			},
			expectedErr: errStore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store := mocks.NewStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(store)
			}

			service := metricservice.New(store)
			err := service.UpdateGauge(ctx, tt.dto)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestService_UpdateCounter(t *testing.T) {
	errStore := errors.New("store error")
	tests := []struct {
		name        string
		dto         metricservice.UpdateCounterDTO
		mockSetup   func(store *mocks.Store)
		expectedErr error
	}{
		{
			name: "Success update",
			dto: metricservice.UpdateCounterDTO{
				Name:  "test-counter",
				Value: 123,
			},
			mockSetup: func(store *mocks.Store) {
				store.On("UpsertCounterMetric", mock.Anything, metricservice.UpsertCounterMetricParams{
					Name:  "test-counter",
					Value: 123,
				}).Return(true, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Store error",
			dto: metricservice.UpdateCounterDTO{
				Name:  "test-counter",
				Value: 123,
			},
			mockSetup: func(store *mocks.Store) {
				store.On("UpsertCounterMetric", mock.Anything,
					metricservice.UpsertCounterMetricParams{
						Name:  "test-counter",
						Value: 123,
					}).Return(false, errStore)
			},
			expectedErr: errStore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store := mocks.NewStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(store)
			}

			service := metricservice.New(store)
			err := service.UpdateCounter(ctx, tt.dto)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestService_FindAll(t *testing.T) {
	errStore := errors.New("store error")
	tests := []struct {
		name        string
		dto         metricservice.FindAllDTO
		mockSetup   func(store *mocks.Store)
		expectedRes []entity.Metrics
		expectedErr error
	}{
		{
			name: "Success find all metrics",
			dto:  metricservice.FindAllDTO{},
			mockSetup: func(store *mocks.Store) {
				store.On("FindAllMetrics", mock.Anything, metricservice.FindAllMetricsNewParams{}).
					Return([]entity.Metrics{
						{
							ID:    "metric1",
							MType: enum.CounterMetricType,
							Delta: func(v int64) *int64 { return &v }(123),
						},
						{
							ID:    "metric2",
							MType: enum.GaugeMetricType,
							Value: func(v float64) *float64 { return &v }(45.67),
						},
					}, nil)
			},
			expectedRes: []entity.Metrics{
				{
					ID:    "metric1",
					MType: enum.CounterMetricType,
					Delta: func(v int64) *int64 { return &v }(123),
				},
				{
					ID:    "metric2",
					MType: enum.GaugeMetricType,
					Value: func(v float64) *float64 { return &v }(45.67),
				},
			},
			expectedErr: nil,
		},
		{
			name: "Store error",
			dto:  metricservice.FindAllDTO{},
			mockSetup: func(store *mocks.Store) {
				store.On("FindAllMetrics", mock.Anything, metricservice.FindAllMetricsNewParams{}).
					Return(nil, errStore)
			},
			expectedRes: nil,
			expectedErr: errStore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store := mocks.NewStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(store)
			}

			service := metricservice.New(store)
			res, err := service.FindAll(ctx, tt.dto)

			assert.Equal(t, tt.expectedRes, res)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestService_Find(t *testing.T) {
	errStore := errors.New("store error")
	tests := []struct {
		name        string
		dto         metricservice.FindDTO
		mockSetup   func(store *mocks.Store)
		expectedRes entity.Metrics
		expectedErr error
	}{
		{
			name: "Success find gauge metric",
			dto: metricservice.FindDTO{
				MetricType: enum.GaugeMetricType,
				MetricName: "test-gauge",
			},
			mockSetup: func(store *mocks.Store) {
				store.On("FindGaugeMetrics", mock.Anything, metricservice.FindGaugeMetricParams{
					MetricName: "test-gauge",
				}).Return(entity.Metrics{
					ID:    "test-gauge",
					MType: enum.GaugeMetricType,
					Value: func(v float64) *float64 { return &v }(42.0),
				}, nil)
			},
			expectedRes: entity.Metrics{
				ID:    "test-gauge",
				MType: enum.GaugeMetricType,
				Value: func(v float64) *float64 { return &v }(42.0),
			},
			expectedErr: nil,
		},
		{
			name: "Success find counter metric",
			dto: metricservice.FindDTO{
				MetricType: enum.CounterMetricType,
				MetricName: "test-counter",
			},
			mockSetup: func(store *mocks.Store) {
				store.On("FindCounterMetrics", mock.Anything, metricservice.FindCounterMetricParams{
					MetricName: "test-counter",
				}).Return(entity.Metrics{
					ID:    "test-counter",
					MType: enum.CounterMetricType,
					Delta: func(v int64) *int64 { return &v }(123),
				}, nil)
			},
			expectedRes: entity.Metrics{
				ID:    "test-counter",
				MType: enum.CounterMetricType,
				Delta: func(v int64) *int64 { return &v }(123),
			},
			expectedErr: nil,
		},
		{
			name: "Store error on gauge metric",
			dto: metricservice.FindDTO{
				MetricType: enum.GaugeMetricType,
				MetricName: "test-gauge",
			},
			mockSetup: func(store *mocks.Store) {
				store.On("FindGaugeMetrics", mock.Anything, metricservice.FindGaugeMetricParams{
					MetricName: "test-gauge",
				}).Return(entity.Metrics{}, errStore)
			},
			expectedRes: entity.Metrics{},
			expectedErr: errStore,
		},
		{
			name: "Undefined metric type",
			dto: metricservice.FindDTO{
				MetricType: "unknown",
				MetricName: "test-unknown",
			},
			mockSetup:   nil,
			expectedRes: entity.Metrics{},
			expectedErr: errorz.ErrUndefinedMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store := mocks.NewStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(store)
			}

			service := metricservice.New(store)
			res, err := service.Find(ctx, tt.dto)

			if tt.expectedErr != nil {
				require.Error(t, err, "expected an error but got nil")
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err, "did not expect an error but got one")
				assert.Equal(t, tt.expectedRes, res, "result did not match expected")
			}
		})
	}
}

func TestService_UpdateBatch(t *testing.T) {
	errStore := errors.New("store error")
	tests := []struct {
		name        string
		dto         metricservice.UpdateBatchDTO
		mockSetup   func(store *mocks.Store)
		expectedErr error
	}{
		{
			name: "Success update batch",
			dto: metricservice.UpdateBatchDTO{
				Data: &[]entity.Metrics{
					{
						ID:    "metric1",
						MType: enum.CounterMetricType,
						Delta: func(v int64) *int64 { return &v }(123),
					},
					{
						ID:    "metric2",
						MType: enum.GaugeMetricType,
						Value: func(v float64) *float64 { return &v }(45.67),
					},
				},
			},
			mockSetup: func(store *mocks.Store) {
				store.On("UpdateMetricsBatchTx", mock.Anything, metricservice.UpdateMetricsBatchTxParams{
					Data: &[]entity.Metrics{
						{
							ID:    "metric1",
							MType: enum.CounterMetricType,
							Delta: func(v int64) *int64 { return &v }(123),
						},
						{
							ID:    "metric2",
							MType: enum.GaugeMetricType,
							Value: func(v float64) *float64 { return &v }(45.67),
						},
					},
				}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Store error",
			dto: metricservice.UpdateBatchDTO{
				Data: &[]entity.Metrics{
					{
						ID:    "metric1",
						MType: enum.CounterMetricType,
						Delta: func(v int64) *int64 { return &v }(123),
					},
					{
						ID:    "metric2",
						MType: enum.GaugeMetricType,
						Value: func(v float64) *float64 { return &v }(45.67),
					},
				},
			},
			mockSetup: func(store *mocks.Store) {
				store.On("UpdateMetricsBatchTx", mock.Anything, metricservice.UpdateMetricsBatchTxParams{
					Data: &[]entity.Metrics{
						{
							ID:    "metric1",
							MType: enum.CounterMetricType,
							Delta: func(v int64) *int64 { return &v }(123),
						},
						{
							ID:    "metric2",
							MType: enum.GaugeMetricType,
							Value: func(v float64) *float64 { return &v }(45.67),
						},
					},
				}).Return(errStore)
			},
			expectedErr: errStore,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			store := mocks.NewStore(t)
			if tt.mockSetup != nil {
				tt.mockSetup(store)
			}

			service := metricservice.New(store)
			err := service.UpdateBatch(ctx, tt.dto)

			if tt.expectedErr != nil {
				require.Error(t, err, "expected an error but got nil")
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err, "did not expect an error but got one")
			}
		})
	}
}
