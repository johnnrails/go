package routers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/repository"
	"github.com/johnnrails/ddd_go/with_gin/users/serializers"
)

func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/:username", ProfileRetrieve)
	router.POST("/:username/follow", ProfileFollow)
	router.DELETE("/:username/follow", ProfileUnfollow)
}

func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")

	repo := repository.UserRepository{
		DB: common.GetDB(),
	}

	userModel, err := repo.FindOneUser(&models.UserModel{
		Username: username,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}

	serializer := serializers.ProfileSerializer{userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response(c)})
}

func ProfileFollow(c *gin.Context) {
	username := c.Param("username")

	repo := repository.UserRepository{
		DB: common.GetDB(),
	}

	userModel, err := repo.FindOneUser(&models.UserModel{
		Username: username,
	})

	myUserModel := c.MustGet("user_model").(models.UserModel)

	err = repo.Following(myUserModel, userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := serializers.ProfileSerializer{userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response(c)})
}

func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")

	repo := repository.UserRepository{
		DB: common.GetDB(),
	}

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

	ctx := context.WithValue(c, "user_model", model)
	serializer := serializers.ProfileSerializer{}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response(ctx)})
}
