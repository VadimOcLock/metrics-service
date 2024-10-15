package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func healthHandler(pool *pgx.Conn) http.Handler {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
