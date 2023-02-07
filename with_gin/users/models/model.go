package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/johnnrails/ddd_go/with_gin/common"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email;unique_index"`
	Bio      string `gorm:"column:bio;size:1024"`
	Password string `gorm:"column:password;not null"`
}

type FollowModel struct {
	gorm.Model
	Following    UserModel
	FollowingID  uint
	FollowedBy   UserModel
	FollowedByID uint
}

func AutoMigrate() {
	db := common.GetDB()
	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&FollowModel{})
}

func (u *UserModel) SetPassword(pass string) error {
	if len(pass) < 8 {
		return errors.New("Password should have more than 8 characters")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hash)
	return nil
}

func (u *UserModel) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
