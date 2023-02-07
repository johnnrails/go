package users

import (
	"context"
	"strings"

	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

// Strips 'TOKEN ' prefix from token string
func GetTokenFromBearerToken(tok string) (string, error) {
	// Should be a bearer token
	if strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

func UpdateContextUserModel(ctx context.Context, id uint) context.Context {
	var model models.UserModel
	if id != 0 {
		db := common.GetDB()
		db.First(&model, id)
	}
	c := context.WithValue(ctx, "user_model", model)
	return c
}
