package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
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

		wrappedWriter := NewResponseWriterWrapper(w)

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)

		log.Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Dur("duration", duration).
			Int("status", wrappedWriter.statusCode).
			Int("content length", wrappedWriter.contentLength).
			Msg("request completed")
	})
}

// GZIP middleware обрабатывает сжатие и декомпрессию данных.
func GZIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что клиент поддерживает получение сжатых данных в формате gzip.
		if supportsGZIP(r) {
			cw := newCompressWriter(w)
			defer func() {
				if err := cw.Close(); err != nil {
					log.Error().Err(err).Send()
				}
			}()
			w = cw
		}

		// Проверяем, что клиент отправил серверу сжатые данные в формате gzip.
		if sendsGZIP(r) {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusInternalServerError)

				return
			}
			defer func() {
				if err = cr.Close(); err != nil {
					log.Error().Err(err).Send()
				}
			}()
			r.Body = cr
		}

		next.ServeHTTP(w, r)
	})
}

// supportsGZIP проверяет, поддерживает ли клиент получение сжатых данных в формате gzip от сервера.
func supportsGZIP(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

// sendsGZIP проверяет, отправил ли клиент серверу сжатые данные в формате gzip.
func sendsGZIP(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки.
type compressWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	gw := gzip.NewWriter(w)
	return &compressWriter{
		ResponseWriter: w,
		gzipWriter:     gw,
	}
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	return cw.gzipWriter.Write(b)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		cw.Header().Set("Content-Encoding", "gzip")
	}
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.gzipWriter.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные.
type compressReader struct {
	io.ReadCloser
	gzipReader *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		ReadCloser: r,
		gzipReader: gr,
	}, nil
}

func (cr *compressReader) Read(p []byte) (int, error) {
	return cr.gzipReader.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.ReadCloser.Close(); err != nil {
		return err
	}
	return cr.gzipReader.Close()
}
