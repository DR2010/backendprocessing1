package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.SetOutput(os.Stderr)

		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.Host,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}