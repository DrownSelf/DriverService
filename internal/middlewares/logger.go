package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s|%v|%s|%s|%s",
			r.Method,
			r.Response.StatusCode,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
