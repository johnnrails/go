package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/johnnrails/ddd_go/boilerplate/logger"
	"github.com/johnnrails/ddd_go/boilerplate/metadata"
	"github.com/vardius/gorouter/v4"
)

func Logger() gorouter.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			logger.Info(r.Context(), fmt.Sprintf("[HTTP START]: %s %s %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			))

			next.ServeHTTP(w, r)
			end := time.Since(now)

			statusCode := http.StatusOK
			stackTrace := ""

			if m := metadata.FromContext(r.Context()); m != nil {
				statusCode = m.StatusCode
				if m.Err != nil {
					stackTrace = m.Err.Error()
				}
			}

			args := []interface{}{
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				statusCode,
				end,
			}

			msg := fmt.Sprintf("[HTTP] End: %s %s -> %s [%d] (%s): %s", args...)

			if stackTrace == "" {
				logger.Info(r.Context(), msg)
			}

			args = append(args, stackTrace)
			if statusCode == http.StatusInternalServerError {
				logger.Error(r.Context(), msg)
			} else {
				logger.Debug(r.Context(), msg)
			}
		})
	}
}
