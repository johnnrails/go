package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type CommentValidator struct {
	Comment struct {
		Body string `form:"body" json:"body" binding:"max=2048"`
	} `json:"comment"`
	commentModel Comment `json:"-"`
}

func NewCommentValidator() CommentValidator {
	return CommentValidator{}
}

func (v *CommentValidator) Bind(c *gin.Context) error {
	userModel := c.MustGet("user_model").(models.UserModel)

	err := common.Bind(c, v)
	if err != nil {
		return err
	}

	v.commentModel.Body = v.Comment.Body
	v.commentModel.AuthorID = AuthorRepository{
		DB: common.GetDB(),
	}.GetAuthor(userModel).ID

	return nil
}
