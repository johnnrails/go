package auth

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	r := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Host: "localhost:8000",
		},
		Header: http.Header{
			"Authorization": []string{"BEARER eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2NzYxMjIzODksImV4cCI6MTcwNzY1ODM4OSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIm5hbWUiOiJKb2hubnkifQ.03bjc7COQLjNguIdG3jl2SkhNljEjUeYem9uoPLa1QI"},
		},
	}
	var claims jwt.MapClaims
	token, _ := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		// This is the lookup key function, use this if you want to verify yourself the token
		func(t *jwt.Token) (i interface{}, e error) {
			log.Println(i)
			log.Println(t)
			// Return the secret key if the token is valid
			return []byte("qwertyuiopasdfghjklzxcvbnm123456"), nil
		},
		request.WithClaims(&claims),
	)

	log.Println(token)
	log.Println(claims)
	assert.Equal(t, token.Valid, true)
}
