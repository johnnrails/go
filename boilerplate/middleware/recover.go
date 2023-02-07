package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	apperrors "github.com/johnnrails/ddd_go/boilerplate/errors"
	"github.com/johnnrails/ddd_go/boilerplate/http/response/json"
	"github.com/johnnrails/ddd_go/boilerplate/logger"
	"github.com/vardius/gorouter/v4"
)

func Recover() gorouter.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Critical(r.Context(), fmt.Sprintf("[HTTP REQUEST RECOVERED]: %v %s", rec, debug.Stack()))
					apperr := apperrors.NewAppErrorFromError(
						fmt.Errorf("%w: recovered from panic", apperrors.ErrInternal),
					)
					if err := json.JSONError(r.Context(), w, http.StatusInternalServerError, apperr); err != nil {
						logger.Critical(r.Context(), fmt.Sprintf("[HTTP ERROR]: While sending response: %v", err))
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
