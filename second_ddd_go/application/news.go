package application

import (
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain/repositories"
)

type NewsApplication struct {
	repository repositories.NewsRepository
}

func CreateNewsApplication(nr repositories.NewsRepository) *NewsApplication {
	return &NewsApplication{
		repository: nr,
	}
}

func (na *NewsApplication) GetByID(id int) (*domain.News, error) {
	return na.repository.Get(id)
}

func (na *NewsApplication) GetAll(limit int, page int) ([]domain.News, error) {
	return na.repository.GetAll()
}

func (na *NewsApplication) Add(n domain.News) error {
	return na.repository.Save(&n)
}

func (na *NewsApplication) Update(id int, new domain.News) error {
	return na.repository.Update(id, new)
}

func (na *NewsApplication) Remove(id int) error {
	return na.repository.Remove(id)
}

func (na *NewsApplication) GetByStatus(status string) ([]domain.News, error) {
	return na.repository.GetByStatus(status)
}

func (na *NewsApplication) GetBySlug(slug string) ([]*domain.News, error) {
	return na.repository.GetBySlug(slug)
}
