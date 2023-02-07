package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/johnnrails/ddd_go/third_ddd_go/common/errors/http_errors"
)

const (
	userContextKey int = iota
)

type User struct {
	UUID  string
	Email string
	Role  string

	DisplayName string
}

func HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims jwt.MapClaims
		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			func(t *jwt.Token) (i interface{}, e error) {
				return []byte("mock_secret"), nil
			},
			request.WithClaims(&claims),
		)

		if err != nil {
			http_errors.BadRequest("unable-to-get-jwt", w, r)
			return
		}

		if !token.Valid {
			http_errors.BadRequest("invalid-jwt", w, r)
			return
		}

		r = r.WithContext(context.WithValue(
			r.Context(),
			userContextKey,
			User{
				UUID:        claims["user_uuid"].(string),
				Email:       claims["email"].(string),
				Role:        claims["role"].(string),
				DisplayName: claims["name"].(string),
			},
		))

		next.ServeHTTP(w, r)
	})
}

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}
	return u, nil
}
