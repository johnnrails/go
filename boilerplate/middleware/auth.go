package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	apperrors "github.com/johnnrails/ddd_go/boilerplate/errors"
	"github.com/johnnrails/ddd_go/boilerplate/identity"
	"github.com/johnnrails/ddd_go/boilerplate/logger"
)

// GetIdentityFromTokenAuthFunc returns Identity from token
type GetIdentityFromTokenAuthFunc func(ctx context.Context, token string) (*identity.Identity, error)

type tokenAuth struct {
	fn GetIdentityFromTokenAuthFunc
}

func (a *tokenAuth) FromHeader(realm string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}

			i, err := a.fn(r.Context(), token[7:])
			if err != nil {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, realm))
				logger.Warning(r.Context(), fmt.Sprintf("[HTTP] failed to authenticate from header: %v", apperrors.NewAppErrorFromError(err)))
			}

			ctx := identity.ContextWithIdentity(r.Context(), i)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *tokenAuth) FromQuery(name string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get(name)

			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			i, err := a.fn(r.Context(), token)
			if err != nil {
				logger.Warning(r.Context(), fmt.Sprintf("[HTTP] failed to authenticate from query: %v", apperrors.NewAppErrorFromError(err)))
			}

			ctx := identity.ContextWithIdentity(r.Context(), i)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *tokenAuth) FromCookie(name string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(name)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			i, err := a.fn(r.Context(), cookie.Value)
			if err != nil {
				logger.Warning(r.Context(), fmt.Sprintf("[HTTP] failed to authenticate from cookie: %v", apperrors.NewAppErrorFromError(err)))
			}

			ctx := identity.ContextWithIdentity(r.Context(), i)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewToken returns new token authenticator
func NewToken(fn GetIdentityFromTokenAuthFunc) *tokenAuth {
	return &tokenAuth{
		fn: fn,
	}
}

// Credential

type GetIdentityFromUserAndPasswordAuthFunc func(username, password string) (identity.Identity, error)

type credentialsAuth struct {
	fn GetIdentityFromUserAndPasswordAuthFunc
}

func (ca *credentialsAuth) FromBasicAuth(realm string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			usr, pass, hasAuth := r.BasicAuth()

			if !hasAuth {
				next.ServeHTTP(w, r)
				return
			}

			i, err := ca.fn(usr, pass)
			if err != nil {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
				logger.Warning(r.Context(), fmt.Sprintf("[HTTP] basic auth failed: %v", apperrors.NewAppErrorFromError(err)))
			}

			ctx := identity.ContextWithIdentity(r.Context(), &i)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewCredentials returns new credentials authenticator
func NewCredentials(fn GetIdentityFromUserAndPasswordAuthFunc) *credentialsAuth {
	return &credentialsAuth{
		fn: fn,
	}
}
