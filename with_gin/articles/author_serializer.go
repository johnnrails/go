package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/users/serializers"
)

type AuthorSerializer struct {
	Author
}

func (s *AuthorSerializer) Response(c *gin.Context) serializers.ProfileResponse {
	response := serializers.ProfileSerializer{s.Author.User}
	return response.Response(c)
}
