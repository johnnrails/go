package articles

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
)

func ArticlesRouter(router *gin.RouterGroup) {
	db := common.GetDB()

	generalRepository := GeneralOperationsRepository{
		DB: db,
	}

	articleRepository := ArticleRepository{
		DB: db,
	}

	router.POST("/", func(ctx *gin.Context) {
		validator := NewArticleValidator()

		if err := validator.Bind(ctx); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
			return
		}

		if err := generalRepository.Insert(&validator.articleModel); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}

		serializer := ArticleSerializer{validator.articleModel}
		ctx.JSON(http.StatusCreated, gin.H{"article": serializer.Response(ctx)})
	})

	router.PUT("/:slug", func(ctx *gin.Context) {
		slug := ctx.Param("slug")
		article, err := articleRepository.FindOne(&Article{
			Slug: slug,
		})

		if err != nil {
			ctx.JSON(http.StatusNotFound, common.NewError("article", errors.New("Invalid Slug")))
			return
		}

		validator := NewArticleValidatorFullWith(article)
		if err := validator.Bind(ctx); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
			return
		}

		validator.articleModel.ID = article.ID

		if err := generalRepository.Update(article, validator.articleModel); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
			return
		}

		serializer := ArticleSerializer{article}
		ctx.JSON(http.StatusOK, gin.H{"article": serializer.Response(ctx)})
	})
}
