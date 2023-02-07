package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users"
)

var AuthorizationFromHeaderExtractor = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	users.GetTokenFromBearerToken,
}

var AuthorizationOrAccessTokenExtractor = &request.MultiExtractor{
	AuthorizationFromHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		users.UpdateContextUserModel(c, 0)

		token, err := request.ParseFromRequest(c.Request, AuthorizationOrAccessTokenExtractor, func(t *jwt.Token) (interface{}, error) {
			return []byte(common.NBSecretPassword), nil
		})

		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := uint(claims["id"].(float64))
			users.UpdateContextUserModel(c, id)
		}
	}
}
