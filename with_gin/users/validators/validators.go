package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type UserFromUserValidator struct {
	Username string `form:"username" json:"username" binding:"exists,alphanum,min=8,max=32"`
	Email    string `form:"email" json:"email" binding:"exists,email"`
	Password string `form:"password" json:"password" binding:"exists,min=12,max=64"`
	Bio      string `form:"bio" json:"bio" binding:"max=1024"`
}

type UserValidator struct {
	UserFromUserValidator `json:"user"`
	UserModel             models.UserModel `json:"-"`
}

func (self *UserValidator) BindFromContext(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.UserModel.Username = self.UserFromUserValidator.Username
	self.UserModel.Email = self.UserFromUserValidator.Email
	self.UserModel.Bio = self.UserFromUserValidator.Bio
	if self.UserFromUserValidator.Password != common.NBRandomPassword {
		self.UserModel.SetPassword(self.UserFromUserValidator.Password)
	}
	return nil
}

func (self *UserValidator) BindFromModel(userModel models.UserModel) {
	self.UserFromUserValidator = UserFromUserValidator{
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password,
		Bio:      userModel.Bio,
	}
}

type LoginValidator struct {
	User struct {
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password" json:"password" binding:"exists,min=12,max=64"`
	}
	UserModel models.UserModel `json:"-"`
}

func (self *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.UserModel.Email = self.User.Email
	return nil
}
