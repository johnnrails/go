package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/serializers"
)

type ArticleSerializer struct {
	Article
}

type ArticleResponse struct {
	ID             uint                        `json:"-"`
	Title          string                      `json:"title"`
	Slug           string                      `json:"slug"`
	Description    string                      `json:"description"`
	Body           string                      `json:"body"`
	CreatedAt      string                      `json:"createdAt"`
	UpdatedAt      string                      `json:"updatedAt"`
	Author         serializers.ProfileResponse `json:"author"`
	Tags           []string                    `json:"tagList"`
	Favorite       bool                        `json:"favorited"`
	FavoritesCount uint                        `json:"favoritesCount"`
}

func (s *ArticleSerializer) Response(c *gin.Context) ArticleResponse {
	userModel := c.MustGet("user_model").(models.UserModel)
	serializer := AuthorSerializer{s.Author}

	repository := ArticleRepository{
		DB: common.GetDB(),
	}

	response := ArticleResponse{
		ID:             s.ID,
		Slug:           slug.Make(s.Title),
		Title:          s.Title,
		Description:    s.Description,
		Body:           s.Body,
		CreatedAt:      s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt:      s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:         serializer.Response(c),
		Favorite:       repository.IsArticleFavoriteBy(s.ID, userModel.ID),
		FavoritesCount: repository.FavoritesCount(s.ID),
	}

	tagSerializer := TagsSerializer{Tags: s.Tags}
	response.Tags = tagSerializer.Response()

	return response
}

type ArticlesSerializer struct {
	Articles []Article
}

func (s *ArticlesSerializer) Response(c *gin.Context) []ArticleResponse {
	response := []ArticleResponse{}
	for _, article := range s.Articles {
		serializer := ArticleSerializer{article}
		response = append(response, serializer.Response(c))
	}
	return response
}
