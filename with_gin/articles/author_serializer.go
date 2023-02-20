package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/users/response"
)

type AuthorSerializer struct {
	Author
}

func (s *AuthorSerializer) Response(c *gin.Context) response.UserResponse {
	return response.ToUserResponse(s.Author.User)
}
