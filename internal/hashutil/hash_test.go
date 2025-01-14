package hashutil_test

import (
	"strconv"
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/hashutil"
)

func TestComputeHMAC(t *testing.T) {
	tests := []struct {
		name     string
		message  []byte
		key      string
		expected string // Ожидаемый результат
	}{
		{
			name:     "Empty message and key",
			message:  []byte{},
			key:      "",
			expected: "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
		},
		{
			name:     "Simple message and key",
			message:  []byte("Hello, World!"),
			key:      "secret",
			expected: "fcfaffa7fef86515c7beb6b62d779fa4ccf092f2e61c164376054271252821ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashutil.ComputeHMAC(tt.message, tt.key)
			if result != tt.expected {
				t.Errorf("ComputeHMAC() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkComputeHMAC(b *testing.B) {
	dataSizes := []int{32, 128, 1024, 10 * 1024} // 32B, 128B, 1KB, 10KB
	key := "benchmark-key"

	for _, size := range dataSizes {
		b.Run(getSizeLabel(size), func(b *testing.B) {
			message := make([]byte, size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashutil.ComputeHMAC(message, key)
			}
		})
	}
}

func getSizeLabel(size int) string {
	if size >= 1024 {
		return "Size-" + strconv.Itoa(size/1024) + "KB"
	}

	return "Size-" + strconv.Itoa(size) + "B"
}
