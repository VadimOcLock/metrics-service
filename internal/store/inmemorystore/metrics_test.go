package inmemorystore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/VadimOcLock/metrics-service/internal/service/metricservice"
	"github.com/VadimOcLock/metrics-service/internal/store/inmemorystore"
	"github.com/stretchr/testify/assert"
)

func TestImpl_UpsertGaugeMetric(t *testing.T) {
	t.Run("Insert new gauge metric", func(t *testing.T) {
		store := inmemorystore.New()

		arg := metricservice.UpsertGaugeMetricParams{
			Name:  "metric1",
			Value: 42.0,
		}
		ok, err := store.UpsertGaugeMetric(context.Background(), arg)

		require.NoError(t, err)
		assert.True(t, ok)

		store.S.Mu.RLock()
		defer store.S.Mu.RUnlock()
		assert.Equal(t, 42.0, store.S.Gauges["metric1"])
	})

	t.Run("Update existing gauge metric", func(t *testing.T) {
		store := inmemorystore.New()

		_, err := store.UpsertGaugeMetric(context.Background(), metricservice.UpsertGaugeMetricParams{
			Name:  "metric1",
			Value: 42.0,
		})
		require.NoError(t, err)

		arg := metricservice.UpsertGaugeMetricParams{
			Name:  "metric1",
			Value: 84.0,
		}
		ok, err := store.UpsertGaugeMetric(context.Background(), arg)

		require.NoError(t, err)
		assert.True(t, ok)

		store.S.Mu.RLock()
		defer store.S.Mu.RUnlock()
		assert.Equal(t, 84.0, store.S.Gauges["metric1"])
	})
}

func BenchmarkImpl_UpsertGaugeMetric_Insert(b *testing.B) {
	store := inmemorystore.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		name := fmt.Sprintf("metric-%d", i)
		arg := metricservice.UpsertGaugeMetricParams{
			Name:  name,
			Value: float64(i),
		}

		b.StartTimer()
		_, err := store.UpsertGaugeMetric(context.Background(), arg)
		if err != nil {
			b.Fatalf("error inserting metric: %v", err)
		}
	}
}

func BenchmarkImpl_UpsertGaugeMetric_Update(b *testing.B) {
	store := inmemorystore.New()

	_, err := store.UpsertGaugeMetric(context.Background(), metricservice.UpsertGaugeMetricParams{
		Name:  "metric1",
		Value: 42.0,
	})
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arg := metricservice.UpsertGaugeMetricParams{
			Name:  "metric1",
			Value: float64(i),
		}
		_, err = store.UpsertGaugeMetric(context.Background(), arg)
		if err != nil {
			b.Fatalf("error updating metric: %v", err)
		}
	}
}
