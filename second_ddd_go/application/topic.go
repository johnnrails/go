package application

import (
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain/repositories"
)

type TopicApplication struct {
	repo repositories.TopicRepository
}

func CreateTopicApplication(tr repositories.TopicRepository) *TopicApplication {
	return &TopicApplication{
		repo: tr,
	}
}

func (ta *TopicApplication) Get(id int) (*domain.Topic, error) {
	return ta.repo.Get(id)
}

func (ta *TopicApplication) GetAll() ([]domain.Topic, error) {
	return ta.repo.GetAll()
}

func (ta *TopicApplication) Add(name string, slug string) error {
	return ta.repo.Save(&domain.Topic{
		Name: name,
		Slug: slug,
	})
}

func (ta *TopicApplication) Remove(id int) error {
	return ta.repo.Remove(id)
}

func (ta *TopicApplication) Update(t domain.Topic, id int) error {
	t.ID = uint(id)
	return ta.repo.Update(&t)
}
