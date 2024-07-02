package handler

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Request: %s %s\n", req.Method, req.URL.Path)
		next.ServeHTTP(res, req)
	})
}
