package middleware

import (
	"fmt"
	apperrors "github.com/johnnrails/ddd_go/boilerplate/errors"
	"github.com/johnnrails/ddd_go/boilerplate/http/response/json"
	"github.com/johnnrails/ddd_go/boilerplate/identity"
	"net/http"
)

// GrantAccessFor returns Status Unauthorized if
// Identity is not set within request's context
// or user does not have required permission
func CheckAccessFor(permission identity.Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			i, ok := identity.FromContext(r.Context())

			if !ok {
				json.JSONError(
					r.Context(),
					w,
					http.StatusUnauthorized,
					apperrors.NewAppErrorFromError(fmt.Errorf("%w: request is missing identity", apperrors.ErrUnauthorized)),
				)
				return
			}

			if !i.Permission.Has(permission) {
				json.JSONError(
					r.Context(),
					w,
					http.StatusForbidden,
					apperrors.NewAppErrorFromError(
						fmt.Errorf("%w: (%d) missing permission %d",
							apperrors.ErrForbidden,
							i.Permission,
							permission,
						),
					),
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
