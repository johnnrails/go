package repositories

import "github.com/johnnrails/ddd_go/second_ddd_go/domain"

type TopicRepository interface {
	Get(id int) (*domain.Topic, error)
	GetAll() ([]domain.Topic, error)
	Save(topic *domain.Topic) error
	Remove(id int) error
	Update(topic *domain.Topic) error
}
