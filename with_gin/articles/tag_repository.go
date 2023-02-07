package articles

import "github.com/jinzhu/gorm"

type TagRepository struct {
	DB *gorm.DB
}

func (tr TagRepository) GetAll() ([]Tag, error) {
	var tags []Tag
	err := tr.DB.Find(&tags).Error
	return tags, err
}
