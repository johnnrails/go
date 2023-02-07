package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type ArticleValidator struct {
	Article struct {
		Title       string   `form:"title" json:"title" binding:"exists,min4"`
		Description string   `form:"description" json:"description" binding:"max=2048"`
		Body        string   `form:"body" json:"body" binding:"max=2048"`
		Tags        []string `form:"tagList" json:"tagList"`
	} `json:"article"`
	articleModel Article `json:"-"`
}

func NewArticleValidator() ArticleValidator {
	return ArticleValidator{}
}

func NewArticleValidatorFullWith(article Article) ArticleValidator {
	validator := NewArticleValidator()
	validator.Article.Title = article.Title
	validator.Article.Description = article.Description
	validator.Article.Body = article.Body

	for _, tagModel := range article.Tags {
		validator.Article.Tags = append(validator.Article.Tags, tagModel.Tag)
	}

	return validator
}

func (v *ArticleValidator) Bind(c *gin.Context) error {
	userModel := c.MustGet("user_model").(models.UserModel)
	err := common.Bind(c, v)

	if err != nil {
		return err
	}

	v.articleModel = Article{
		Slug:        slug.Make(v.Article.Title),
		Title:       v.Article.Title,
		Description: v.Article.Description,
		Body:        v.Article.Body,
		Author: AuthorRepository{
			DB: common.GetDB(),
		}.GetAuthor(userModel),
	}

	ArticleRepository{
		DB: common.GetDB(),
	}.SetTags(v.articleModel, v.Article.Tags)

	return nil
}
