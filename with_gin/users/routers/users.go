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

	repo := repository.UserRepository{DB: common.GetDB()}
	if err := repo.SaveOne(validator.UserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": response.ToUserResponse(validator.UserModel)})
}

func UsersLogin(c *gin.Context) {
	loginValidator := validators.LoginValidator{}

	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	repo := repository.UserRepository{DB: common.GetDB()}
	userModel, err := repo.FindOneUser(&models.UserModel{
		Email: loginValidator.UserModel.Email,
	})

	if err != nil || userModel.CheckPassword(loginValidator.UserModel.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not registered")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": response.ToUserResponse(userModel)})
}
