package articles

import (
	"github.com/jinzhu/gorm"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type ArticleRepository struct {
	DB *gorm.DB
}

func (ar ArticleRepository) FavoritesCount(articleID uint) uint {
	var count uint
	ar.DB.Model(&Favorite{}).Where(Favorite{
		FavoriteID: articleID,
	}).Count(&count)
	return count
}

func (ar ArticleRepository) IsArticleFavoriteBy(articleID uint, userID uint) bool {
	var favorite Favorite
	ar.DB.Model(&Author{
		UserID: userID,
	}).Where(Favorite{
		FavoriteID:   articleID,
		FavoriteByID: userID,
	}).First(&favorite)
	return favorite.ID != 0
}

func (ar ArticleRepository) FavoriteArticle(articleID uint, userID uint) error {
	var favorite Favorite
	err := ar.DB.FirstOrCreate(&favorite, &Favorite{
		FavoriteID:   articleID,
		FavoriteByID: userID,
	}).Error
	return err
}

func (ar ArticleRepository) UnfavoriteArticle(articleID uint, userID uint) error {
	err := ar.DB.Where(&Favorite{
		FavoriteID:   articleID,
		FavoriteByID: userID,
	}).Delete(Favorite{}).Error
	return err
}

func (ar ArticleRepository) FindOne(condition interface{}) (Article, error) {
	var model Article

	tx := ar.DB.Begin()

	tx.Where(condition).First(&model)
	tx.Model(&model).Related(&model.Author, "Author")
	tx.Model(&model.Author).Related(&model.Author.User)
	tx.Model(&model).Related(&model.Tags, "Tags")

	err := tx.Commit().Error

	return model, err
}

func (ar ArticleRepository) GetComments(article Article) error {
	tx := ar.DB.Begin()
	tx.Model(article).Related(&article.Comments, "Commnets")
	err := tx.Commit().Error
	return err
}

func (ar ArticleRepository) FindMany(limit uint, offset uint) ([]Article, error) {
	var articles []Article
	tx := ar.DB.Begin()
	ar.DB.Offset(offset).Limit(limit).Find(&articles)

	for i := range articles {
		a := &articles[i]
		tx.Model(&a).Related(&a.Author, "Author")
		tx.Model(&a).Related(&a.Author.User)
	}

	err := tx.Commit().Error
	return articles, err
}

func (ar ArticleRepository) FindManyByTag(tag string, limit uint, offset uint) ([]Article, error) {
	var articles []Article

	tx := ar.DB.Begin()

	if tag != "" {
		var tagModel Tag
		tx.Where(Tag{Tag: tag}).First(&tagModel)

		if tagModel.ID != 0 {
			tx.Model(&tag).Offset(offset).Limit(limit).Related(&articles, "Articles")
		}
	} else {
		ar.DB.Offset(offset).Limit(limit).Find(&articles)
	}

	for i := range articles {
		a := &articles[i]
		tx.Model(&a).Related(&a.Author, "Author")
		tx.Model(&a).Related(&a.Author.User)
		tx.Model(&a).Related(&a.Tags, "Tags")
	}

	err := tx.Commit().Error
	return articles, err
}

func (ar ArticleRepository) FindManyByAuthorName(authorName string, offset uint, limit uint) ([]Article, error) {
	var articles []Article
	tx := ar.DB.Begin()

	if authorName != "" {
		var user models.UserModel
		ar.DB.Where(models.UserModel{Username: authorName}).First(&user)

		author := AuthorRepository{
			DB: common.GetDB(),
		}.GetAuthor(user)

		if author.ID != 0 {
			tx.Model(&author).Offset(offset).Limit(limit).Related(&articles, "Articles")
		}
	} else {
		ar.DB.Offset(offset).Limit(limit).Related(&articles, "Articles")
	}

	err := tx.Commit().Error
	return articles, err
}

func (ar ArticleRepository) FindManyByFavoritedName(favoritedName string, offset uint, limit uint) ([]Article, error) {
	var user models.UserModel
	var articles []Article

	tx := ar.DB.Begin()

	tx.Where(models.UserModel{Username: favoritedName}).First(&user)
	articleUserModel := AuthorRepository{
		DB: common.GetDB(),
	}.GetAuthor(user)

	if articleUserModel.ID != 0 {
		var favoriteModels []Favorite

		tx.Where(Favorite{
			FavoriteByID: articleUserModel.ID,
		}).Offset(offset).Limit(limit).Find(&favoriteModels)

		for _, favorite := range favoriteModels {
			var a Article
			tx.Model(&favorite).Related(&a, "Favorite")
			articles = append(articles, a)
		}
	}

	err := tx.Commit().Error

	return articles, err
}

func (ar ArticleRepository) SetTags(article Article, tags []string) error {
	var tagList []Tag
	for _, tag := range tags {
		var tagModel Tag
		err := ar.DB.First(&tagModel, Tag{Tag: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	article.Tags = tagList
	return nil
}
