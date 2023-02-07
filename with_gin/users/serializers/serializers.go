package serializers

import (
	"context"

	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/repository"
)

type ProfileSerializer struct {
	models.UserModel
}

type ProfileResponse struct {
	ID        uint   `json:"-"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Following bool   `json:"following"`
}

func (self *ProfileSerializer) Response(c context.Context) ProfileResponse {
	userModel := c.Value("user_model").(models.UserModel)
	ur := repository.UserRepository{
		DB: common.GetDB(),
	}
	user := ProfileResponse{
		ID:        self.ID,
		Username:  self.Username,
		Bio:       self.Bio,
		Following: ur.IsFollowing(userModel, self.UserModel),
	}
	return user
}

type UserSerializer struct{}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Token    string `json:"token"`
}

func (self *UserSerializer) Response(c context.Context) UserResponse {
	um := c.Value("user_model").(models.UserModel)
	u := UserResponse{
		Username: um.Username,
		Email:    um.Email,
		Bio:      um.Bio,
		Token:    common.GenToken(um.ID),
	}
	return u
}
