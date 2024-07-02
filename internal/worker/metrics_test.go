package worker

import (
	"context"
	"github.com/VadimOcLock/metrics-service/internal/entity"
	"github.com/VadimOcLock/metrics-service/internal/entity/enum"
	"github.com/VadimOcLock/metrics-service/internal/errorz"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Не понимаю как можно протестировать это без моков
func Test_sendMetric(t *testing.T) {
	type want struct {
		errWait bool
	}
	tests := []struct {
		name   string
		metric entity.MetricDTO
		want   want
	}{
		{
			name: "Successful request",
			metric: entity.MetricDTO{
				Type:  enum.GaugeMetricType,
				Name:  "Alloc",
				Value: "123",
			},
			want: want{
				errWait: false,
			},
		},
		{
			name: "Invalid URL",
			metric: entity.MetricDTO{
				Type:  enum.GaugeMetricType,
				Name:  "",
				Value: "12345",
			},
			want: want{
				errWait: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, errorz.ErrMsgOnlyPOSTMethodAccept, http.StatusMethodNotAllowed)

					return
				}
				parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
				for _, part := range parts {
					if part == "" {
						http.Error(w, errorz.ErrMsgEmptyMetricParam, http.StatusNotFound)

						return
					}
				}
			}))
			defer server.Close()

			serverAddress := server.URL
			if tt.name == "Invalid URL" {
				serverAddress = "http://invalid-url"
			}

			opts := sendMetricOpts{
				client:        server.Client(),
				serverAddress: serverAddress,
				metric:        tt.metric,
			}

			err := sendMetric(context.Background(), opts)
			if tt.want.errWait {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
