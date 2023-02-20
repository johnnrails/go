package response

import (
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Token    string `json:"token"`
}

func ToUserResponse(um models.UserModel) UserResponse {
	u := UserResponse{
		Username: um.Username,
		Email:    um.Email,
		Bio:      um.Bio,
		Token:    common.GenToken(um.ID),
	}
	return u
}
