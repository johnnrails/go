package users

import (
	"context"

	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

func UpdateContextUserModel(ctx context.Context, id uint) context.Context {
	var model models.UserModel
	if id != 0 {
		db := common.GetDB()
		db.First(&model, id)
	}
	c := context.WithValue(ctx, "user_model", model)
	return c
}
