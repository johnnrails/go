package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/repository"
	"github.com/johnnrails/ddd_go/with_gin/users/serializers"
	"github.com/johnnrails/ddd_go/with_gin/users/validators"
)

func UserRegister(router *gin.RouterGroup) {
	router.GET("/", UserRetrieve)
	router.PUT("/", UserUpdate)
}

func UserRetrieve(c *gin.Context) {
	serializer := serializers.UserSerializer{}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response(c)})
}

func UserUpdate(c *gin.Context) {
	model := c.MustGet("user_model").(models.UserModel)
	validator := validators.NewUserModelValidatorFillWith(model)

	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	validator.UserModel.ID = model.ID

	repo := repository.UserRepository{
		DB: common.GetDB(),
	}

	if err := repo.Update(model, validator.UserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	ctx := users.UpdateContextUserModel(c, model.ID)
	serializer := serializers.UserSerializer{}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response(ctx)})
}
