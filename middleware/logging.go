package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// log.Printf("Started %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("%s %s in %v", r.Method, r.RequestURI, time.Since(startTime))
		// TODO: structured logging and only 1 line with took and error and all req and res param, method name
	})
}
