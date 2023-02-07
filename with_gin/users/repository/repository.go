package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/johnnrails/ddd_go/with_gin/users/models"
)

type UserRepository struct {
	DB *gorm.DB
}

func (ur UserRepository) FindOneUser(condition interface{}) (models.UserModel, error) {
	var model models.UserModel
	err := ur.DB.Where(condition).First(&model).Error
	return model, err
}

func (ur UserRepository) SaveOne(u models.UserModel) error {
	err := ur.DB.Save(u).Error
	return err
}

func (ur UserRepository) Update(u models.UserModel, data interface{}) error {
	err := ur.DB.Model(u).Update(data).Error
	return err
}

func (ur UserRepository) Following(u models.UserModel, v models.UserModel) error {
	var follow models.FollowModel
	err := ur.DB.FirstOrCreate(&follow, &models.FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Error
	return err
}

func (ur UserRepository) IsFollowing(u models.UserModel, v models.UserModel) bool {
	var follow models.FollowModel
	ur.DB.Where(models.FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).First(&follow)
	return follow.ID != 0
}

func (ur UserRepository) UnFollowing(u models.UserModel, v models.UserModel) error {
	err := ur.DB.Where(models.FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Delete(models.FollowModel{}).Error
	return err
}

func (ur UserRepository) GetFollowings(u models.UserModel) []models.UserModel {
	tx := ur.DB.Begin()

	var follows []models.FollowModel
	var followings []models.UserModel
	// Get Follows
	tx.Where(models.FollowModel{
		FollowedByID: u.ID,
	}).Find(&follows)
	// Get each follow and build a models.UserModel
	for _, f := range follows {
		var userModel models.UserModel
		tx.Model(&f).Related(&userModel, "Following")
		followings = append(followings, userModel)
	}

	tx.Commit()
	return followings
}
