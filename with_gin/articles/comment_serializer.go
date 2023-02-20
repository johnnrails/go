package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/response"
)

type CommentResponse struct {
	ID        uint                  `json:"id"`
	Body      string                `json:"body"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
	Author    response.UserResponse `json:"author"`
}

type CommentSerializer struct {
	Comment
}

func (s *CommentSerializer) Response(c *gin.Context) CommentResponse {
	author := AuthorRepository{
		DB: common.GetDB(),
	}.GetAuthorByID(s.AuthorID)

	serializer := AuthorSerializer{Author: author}

	response := CommentResponse{
		ID:        s.ID,
		Body:      s.Body,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:    serializer.Response(c),
	}

	return response
}

type CommentsSerializer struct {
	Comments []Comment
}
