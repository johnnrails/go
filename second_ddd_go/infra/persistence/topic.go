package persistence

import (
	"github.com/johnnrails/ddd_go/second_ddd_go/config"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain/repositories"
	"gorm.io/gorm"
)

type TopicRepositoryImpl struct {
	DB *gorm.DB
}

func CreateTopicRepository() (repositories.TopicRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	return &TopicRepositoryImpl{DB: db}, nil
}

func (r *TopicRepositoryImpl) Get(id int) (*domain.Topic, error) {
	topic := &domain.Topic{}
	if err := r.DB.Preload("News").First(&topic, id).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

func (r *TopicRepositoryImpl) GetAll() ([]domain.Topic, error) {
	topics := []domain.Topic{}
	if err := r.DB.Preload("News").Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (r *TopicRepositoryImpl) Save(topic *domain.Topic) error {
	if err := r.DB.Save(&topic).Error; err != nil {
		return err
	}

	return nil
}

func (r *TopicRepositoryImpl) Remove(id int) error {
	topic := &domain.Topic{}
	if err := r.DB.Delete(&topic).Error; err != nil {
		return err
	}
	return nil
}

func (r *TopicRepositoryImpl) Update(topic *domain.Topic) error {
	d := &domain.Topic{
		Name: topic.Name,
		Slug: topic.Slug,
	}
	if err := r.DB.UpdateColumns(d).Error; err != nil {
		return err
	}
	return nil
}
