package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/VadimOcLock/metrics-service/internal/hashutil"
)

func RequestSignatureVerificationMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)

				return
			}
			clientHash := r.Header.Get("HashSHA256")
			if clientHash == "" {
				next.ServeHTTP(w, r)

				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)

				return
			}

			serverHash := hashutil.ComputeHMAC(body, key)
			if clientHash != serverHash {
				http.Error(w, "Invalid hash", http.StatusBadRequest)

				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			next.ServeHTTP(w, r)
		})
	}
}

func ResponseSigningMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signedWriter := &SignedResponseWriter{ResponseWriter: w, key: key}
			next.ServeHTTP(signedWriter, r)
		})
	}
}

type SignedResponseWriter struct {
	http.ResponseWriter
	key    string
	buffer bytes.Buffer
}

func (w *SignedResponseWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	bs, err := w.ResponseWriter.Write(data)
	if err != nil {
		return 0, fmt.Errorf("SignedResponseWriter err: %w", err)
	}

	return bs, nil
}

func (w *SignedResponseWriter) WriteHeader(statusCode int) {
	if w.key != "" {
		responseHash := hashutil.ComputeHMAC(w.buffer.Bytes(), w.key)
		w.Header().Set("HashSHA256", responseHash)
	}
	w.ResponseWriter.WriteHeader(statusCode)
}
