package middleware

import (
	"net/http"

	"github.com/vardius/gorouter/v4"
)

// LimitRequestBody limits the request body
func LimitInBytesRequestBody(n int64) gorouter.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}
