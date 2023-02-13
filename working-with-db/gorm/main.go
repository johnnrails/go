package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code       string
	Price      uint
	Categories []*Category `gorm:"many2many:product_categories;"`
}

type Category struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email;unique_index"`
	Bio      string `gorm:"column:bio;size:1024"`
	Password string `gorm:"column:password;not null"`
}

type Follow struct {
	gorm.Model
	FollowedID  uint
	Followed    User
	FollowingID uint
	Following   User
}

type Favorite struct {
	gorm.Model
	ArticleID uint
	Article   Article
	UserID    uint
}

type Article struct {
	gorm.Model
	AuthorID    uint
	Author      User
	Slug        string `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
}

type Tag struct {
	gorm.Model
	Tag      string    `gorm:"unique_index"`
	Articles []Article `gorm:"many2many:article_tags;"`
}

type Comment struct {
	gorm.Model
	Article   Article
	ArticleID uint
	Author    User
	AuthorID  uint
	Body      string `gorm:"size:2048"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Favorite{})
	db.AutoMigrate(&Follow{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Article{})
}
