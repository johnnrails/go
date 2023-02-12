package routers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/repository"
	"github.com/johnnrails/ddd_go/with_gin/users/serializers"
	"github.com/johnnrails/ddd_go/with_gin/users/validators"
)

func UsersRouter(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UsersRegistration(c *gin.Context) {
	validator := validators.UserValidator{}
	if err := validator.BindFromContext(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	repository := repository.UserRepository{
		DB: common.GetDB(),
	}

	if err := repository.SaveOne(validator.UserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	serializer := serializers.UserSerializer{}
	ctx := context.WithValue(c, "user_model", validator.UserModel)
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response(ctx)})
}

func UsersLogin(c *gin.Context) {
	loginValidator := validators.LoginValidator{}

	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	repository := repository.UserRepository{
		DB: common.GetDB(),
	}

	userModel, err := repository.FindOneUser(&models.UserModel{
		Email: loginValidator.UserModel.Email,
	})

	if err != nil || userModel.CheckPassword(loginValidator.UserModel.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not registered")))
		return
	}

	ctx := users.UpdateContextUserModel(c, userModel.ID)
	serializer := serializers.UserSerializer{}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response(ctx)})
}
