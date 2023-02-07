package middleware

import (
	"net/http"

	"github.com/johnnrails/ddd_go/boilerplate/identity"
	"github.com/rs/cors"
	"github.com/vardius/gorouter/v4"
)

var (
	allowedMethods = []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
	allowedHeaders = []string{"*"}
)

func CORS(allowedOrigins []string, debug bool) gorouter.MiddlewareFunc {
	defaultCors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		Debug:            debug,
	})
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if i, isAuthorized := identity.FromContext(r.Context()); isAuthorized && i.ClientDomain != "" {
				cors := cors.New(cors.Options{
					AllowCredentials: true,
					AllowedOrigins:   []string{i.ClientDomain},
					AllowedMethods:   allowedMethods,
					AllowedHeaders:   allowedHeaders,
					Debug:            debug,
				})
				cors.Handler(next).ServeHTTP(w, r)
				return
			}
			defaultCors.Handler(next).ServeHTTP(w, r)
		})
	}
}
