package lifecycle_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/VadimOcLock/metrics-service/pkg/lifecycle"

	"github.com/safeblock-dev/wr/taskgroup"
	"github.com/stretchr/testify/require"
)

// Пример реализации WorkerRunner
type ExampleWorkerImpl struct{}

func (w *ExampleWorkerImpl) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// Рабочая логика воркера
			fmt.Println("worker tick")
		}
	}
}

func ExampleWorker() {
	// Создаем task group
	tg := taskgroup.New()

	// Создаем и добавляем воркер
	worker := &ExampleWorkerImpl{}
	tg.Add(lifecycle.Worker(worker))

	if err := tg.Run(); err != nil {
		// Обработка ошибки
	}

	// Output: демонстрирует graceful shutdown воркера
}

func ExampleHTTPServer() {
	// Создаем task group
	tg := taskgroup.New()

	// Создаем HTTP сервер
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Добавляем сервер в task group
	tg.Add(lifecycle.HTTPServer(server))

	if err := tg.Run(); err != nil {
		// Обработка ошибки
	}

	// Output: демонстрирует graceful shutdown HTTP сервера
}

func TestWorkerGracefulShutdown(t *testing.T) {
	worker := &ExampleWorkerImpl{}
	execute, interrupt := lifecycle.Worker(worker)

	// Тестируем graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	errCh := make(chan error)
	go func() {
		errCh <- execute()
	}()

	// Имитируем прерывание
	interrupt(nil)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("worker didn't shutdown in time")
	}
}

func TestHTTPServerGracefulShutdown(t *testing.T) {
	// Тестовый сервер
	server := &http.Server{
		Addr: "localhost:0", // случайный порт
	}

	execute, interrupt := lifecycle.HTTPServer(server)

	// Запускаем сервер
	errCh := make(chan error)
	go func() {
		errCh <- execute()
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	// Имитируем прерывание
	interrupt(nil)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("server didn't shutdown in time")
	}
}
