package compress_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"strconv"
	"testing"

	"github.com/VadimOcLock/metrics-service/internal/compress"
)

func TestGZipCompress(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   []byte{},
			wantErr: false,
		},
		{
			name:    "Small input",
			input:   []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "Large input",
			input:   bytes.Repeat([]byte("A"), 1024*1024), // 1MB
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := compress.GZipCompress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GZipCompress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				r, rErr := gzip.NewReader(bytes.NewReader(compressed))
				if rErr != nil {
					t.Fatalf("Failed to create gzip reader: %v", rErr)
				}
				defer func(r *gzip.Reader) {
					_ = r.Close()
				}(r)

				decompressed, rErr := io.ReadAll(r)
				if rErr != nil {
					t.Fatalf("Failed to decompress data: %v", rErr)
				}
				if !bytes.Equal(decompressed, tt.input) {
					t.Fatalf("Decompressed data does not match original, got = %s, want = %s", decompressed, tt.input)
				}
			}
		})
	}
}

func BenchmarkGZipCompress(b *testing.B) {
	dataSizes := []int{1024, 10 * 1024, 100 * 1024, 1 * 1024 * 1024}
	for _, size := range dataSizes {
		b.Run(getSizeLabel(size), func(b *testing.B) {
			src := bytes.Repeat([]byte("A"), size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := compress.GZipCompress(src)
				if err != nil {
					b.Fatalf("Failed to compress data: %v", err)
				}
			}
		})
	}
}

func getSizeLabel(size int) string {
	if size >= 1024*1024 {
		return "Size-" + strconv.Itoa(size/(1024*1024)) + "MB"
	}
	if size >= 1024 {
		return "Size-" + strconv.Itoa(size/1024) + "KB"
	}

	return "Size-" + strconv.Itoa(size) + "B"
}
