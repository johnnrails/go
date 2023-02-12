package repositories

import "github.com/johnnrails/ddd_go/second_ddd_go/domain"

type NewsRepository interface {
	GetByID(id int) (*domain.News, error)
	GetAll() ([]domain.News, error)
	GetBySlug(slug string) ([]*domain.News, error)
	GetByStatus(status string) ([]domain.News, error)
	Save(*domain.News) error
	Remove(id int) error
	Update(id int, new domain.News) error
}
