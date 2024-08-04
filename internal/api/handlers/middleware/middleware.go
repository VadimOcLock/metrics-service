package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{w, http.StatusOK, 0}
}

func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	bytesWritten, err := w.ResponseWriter.Write(b)
	w.contentLength += bytesWritten

	return bytesWritten, err
}

// Logger логгирует информацию о входящих запросах и результатах обработки запроса.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Чтение и сохранение тела запроса
		var requestBody bytes.Buffer
		tee := io.TeeReader(r.Body, &requestBody)
		body, err := io.ReadAll(tee)
		if err != nil {
			http.Error(w, "can't read request body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(&requestBody) // Восстановление тела запроса

		wrappedWriter := NewResponseWriterWrapper(w)

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)

		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Dur("duration", duration).
			Int("status", wrappedWriter.statusCode).
			Int("content length", wrappedWriter.contentLength).
			Str("request body", string(body)).
			Msg("request completed")

	})
}
