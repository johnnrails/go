package validators

import (
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type UserModelValidator struct {
	User struct {
		Username string `form:"username" json:"username" binding:"exists,alphanum,min=8,max=32"`
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password" json:"password" binding:"exists,min=12,max=64"`
		Bio      string `form:"bio" json:"bio" binding:"max=1024"`
	} `json:"user"`
	UserModel models.UserModel `json:"-"`
}

func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.UserModel.Username = self.User.Username
	self.UserModel.Email = self.User.Email
	self.UserModel.Bio = self.User.Bio
	if self.User.Password != common.NBRandomPassword {
		self.UserModel.SetPassword(self.User.Password)
	}
	return nil
}

func NewUserModelValidatorFillWith(userModel models.UserModel) UserModelValidator {
	userModelValidator := UserModelValidator{}
	userModelValidator.User.Username = userModel.Username
	userModelValidator.User.Email = userModel.Email
	userModelValidator.User.Bio = userModel.Bio
	userModelValidator.User.Password = common.NBRandomPassword
	return userModelValidator
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
