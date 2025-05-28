// Package lifecycle предоставляет функции для управления жизненным циклом
// долгоживущих компонентов приложения, таких как воркеры и HTTP-серверы.
// Пакет интегрируется с taskgroup для graceful shutdown.
package lifecycle

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/safeblock-dev/wr/taskgroup"
)

// httpServerShutdownTTL - таймаут для graceful shutdown HTTP-сервера
const httpServerShutdownTTL = 10 * time.Second

// HTTPServer создает пару функций ExecuteFn/InterruptFn для taskgroup
// для управления жизненным циклом HTTP-сервера.
//
// Параметры:
//   - server: конфигурированный HTTP-сервер
//
// Возвращает:
//   - execute: функция запуска сервера (ListenAndServe)
//   - interrupt: функция остановки сервера (Shutdown с таймаутом)
//
// Пример использования:
//
//	tg.Add(lifecycle.HTTPServer(myServer))
func HTTPServer(server *http.Server) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	execute := func() error {
		log.Debug().Msgf("HTTP server starting at addr: %s...", server.Addr)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		log.Debug().Msg("HTTP server finished")

		return err
	}
	interrupt := func(_ error) {
		ctx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTTL)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Msgf("shutdown server err: %v", err)
		}
		log.Debug().Msg("HTTP server shutdown complete")
	}

	return execute, interrupt
}
