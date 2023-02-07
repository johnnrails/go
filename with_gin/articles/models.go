package articles

import (
	"github.com/jinzhu/gorm"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type Author struct {
	gorm.Model
	User      models.UserModel
	UserID    uint
	Articles  []Article  `gorm:"ForeignKey:AuthorID"`
	Favorites []Favorite `gorm:"ForeignKey:FavoriteByID"`
}

type Article struct {
	gorm.Model
	AuthorID    uint
	Author      Author
	Slug        string `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	Tags        []Tag
	Comments    []Comment
}

type Tag struct {
	gorm.Model
	Tag      string    `gorm:"unique_index"`
	Articles []Article `gorm:"many2many:article_tags;"`
}

type Comment struct {
	gorm.Model
	ArticleID uint
	AuthorID  uint
	Body      string `gorm:"size:2048"`
}

type Favorite struct {
	gorm.Model
	Favorite     Article
	FavoriteID   uint
	FavoriteByID uint
}
