package config_test

import (
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/config"
)

func TestDatabaseConfig_InMemoryMode(t *testing.T) {
	tests := []struct {
		name     string
		config   config.DatabaseConfig
		expected bool
	}{
		{
			name:     "InMemoryMode with empty DSN",
			config:   config.DatabaseConfig{DSN: ""},
			expected: true,
		},
		{
			name:     "Not InMemoryMode with non-empty DSN",
			config:   config.DatabaseConfig{DSN: "some-dsn"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.InMemoryMode()
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
