package articles

import (
	"github.com/jinzhu/gorm"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
	"github.com/johnnrails/ddd_go/with_gin/users/repository"
)

type AuthorRepository struct {
	DB *gorm.DB
}

func (ar AuthorRepository) GetAuthor(user models.UserModel) Author {
	var author Author

	if user.ID == 0 {
		return author
	}

	ar.DB.Where(&Article{
		AuthorID: user.ID,
	}).First(&author)

	author.User = user

	return author
}

func (ar AuthorRepository) GetAuthorByID(id uint) Author {
	var author Author

	if id == 0 {
		return author
	}

	ar.DB.Where(&Article{
		AuthorID: id,
	}).First(&author)

	return author
}

func (ar AuthorRepository) GetArticleFeed(user models.UserModel, limit, offset uint) ([]Article, error) {
	var articles []Article
	tx := ar.DB.Begin()

	followings := repository.UserRepository{
		DB: common.GetDB(),
	}.GetFollowings(user)

	var authorsIDs []uint

	for _, f := range followings {
		a := ar.GetAuthor(f)
		authorsIDs = append(authorsIDs, a.ID)
	}

	tx.Where("author_id in (?)", authorsIDs).Order("updated_at desc").Offset(offset).Limit(limit).Find(&articles)

	for i := range articles {
		a := &articles[i]
		tx.Model(&a).Related(&a.Author, "Author")
		tx.Model(&a.Author).Related(&a.Author.User)
		tx.Model(&a).Related(&a.Tags, "Tags")
	}

	err := tx.Commit().Error
	return articles, err
}
