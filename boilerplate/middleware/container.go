package middleware

import (
	"context"
	"net/http"

	"github.com/vardius/gocontainer"
	"github.com/vardius/gorouter/v4"
)

// WithContainer wraps http.Handler with a container middleware
func WithContainer(requestContainer gocontainer.Container) gorouter.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "container", requestContainer)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
