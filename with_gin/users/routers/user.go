package routers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/repository"
	response "github.com/johnnrails/ddd_go/with_gin/users/response"
	"github.com/johnnrails/ddd_go/with_gin/users/validators"
)

func UserRegister(router *gin.RouterGroup) {
	router.GET("/", UserRetrieve)
	router.PUT("/", UserUpdate)
	router.GET("/:username", ProfileRetrieve)
	router.POST("/:username/follow", ProfileFollow)
	router.DELETE("/:username/follow", ProfileUnfollow)
}

func GetUserFromContext(c *gin.Context) models.UserModel {
	return c.MustGet("user_model").(models.UserModel)
}

func UserRetrieve(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user": response.ToUserResponse(GetUserFromContext(c))})
}

func UserUpdate(c *gin.Context) {
	model := GetUserFromContext(c)

	validator := validators.UserValidator{}
	validator.BindFromModel(model)

	if err := validator.BindFromContext(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	validator.UserModel.ID = model.ID
	repo := repository.UserRepository{DB: common.GetDB()}
	if err := repo.Update(model.ID, validator.UserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": response.ToUserResponse(validator.UserModel)})
}

func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")

	repo := repository.UserRepository{DB: common.GetDB()}
	userModel, err := repo.FindOneUser(&models.UserModel{
		Username: username,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": response.ToUserResponse(userModel)})
}

func ProfileFollow(c *gin.Context) {
	username := c.Param("username")

	repo := repository.UserRepository{DB: common.GetDB()}
	userModel, err := repo.FindOneUser(&models.UserModel{
		Username: username,
	})

	myUserModel := c.MustGet("user_model").(models.UserModel)

	err = repo.Following(myUserModel, userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": response.ToUserResponse(myUserModel)})
}

func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")

	repo := repository.UserRepository{DB: common.GetDB()}
	model, err := repo.FindOneUser(&models.UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}

	myUserModel := c.MustGet("user_model").(models.UserModel)

	err = repo.UnFollowing(myUserModel, model)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": response.ToUserResponse(myUserModel)})
}
