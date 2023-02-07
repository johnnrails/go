package articles

import (
	"github.com/jinzhu/gorm"
)

type GeneralOperationsRepository struct {
	DB *gorm.DB
}

func (g GeneralOperationsRepository) Insert(data interface{}) error {
	err := g.DB.Save(data).Error
	return err
}

func (g GeneralOperationsRepository) Update(actual interface{}, new interface{}) error {
	err := g.DB.Model(actual).Update(new).Error
	return err
}

func (g GeneralOperationsRepository) Delete(condition interface{}, model interface{}) error {
	err := g.DB.Where(model).Delete(model).Error
	return err
}
