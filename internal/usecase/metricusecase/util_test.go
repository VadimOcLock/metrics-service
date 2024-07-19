package metricusecase_test

import (
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/usecase/metricusecase"

	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_sortMetrics(t *testing.T) {
	tests := []struct {
		name    string
		input   []entity.Metric
		expects []entity.Metric
	}{
		{
			name: "already sorted",
			input: []entity.Metric{
				{Type: "a", Name: "1"},
				{Type: "a", Name: "2"},
				{Type: "b", Name: "3"},
			},
			expects: []entity.Metric{
				{Type: "a", Name: "1"},
				{Type: "a", Name: "2"},
				{Type: "b", Name: "3"},
			},
		},
		{
			name: "unsorted",
			input: []entity.Metric{
				{Type: "b", Name: "3"},
				{Type: "a", Name: "2"},
				{Type: "a", Name: "1"},
			},
			expects: []entity.Metric{
				{Type: "a", Name: "1"},
				{Type: "a", Name: "2"},
				{Type: "b", Name: "3"},
			},
		},
		{
			name: "all same type",
			input: []entity.Metric{
				{Type: "a", Name: "2"},
				{Type: "a", Name: "1"},
				{Type: "a", Name: "3"},
			},
			expects: []entity.Metric{
				{Type: "a", Name: "1"},
				{Type: "a", Name: "2"},
				{Type: "a", Name: "3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]entity.Metric, len(tt.input))
			copy(inputCopy, tt.input)

			metricusecase.SortMetrics(&inputCopy)

			assert.Equal(t, tt.expects, inputCopy)
		})
	}
}
