package middlewares

import (
	"log"
	"net/http"
)

// LogRequest logs the request method, path and IP address
// of the request to the standard output.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request %s %s from IP %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
