package middleware

import (
	"net/http"

	"github.com/vardius/gorouter/v4"
)

// XSS (Cross-site-scripting) sets xss response header types
func XSS() gorouter.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Content-Type-Options", "nosniff")
			w.Header().Add("X-Frame-Options", "DENY")
			next.ServeHTTP(w, r)
		})
	}
}
