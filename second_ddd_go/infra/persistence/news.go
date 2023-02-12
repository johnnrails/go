package persistence

import (
	"github.com/johnnrails/ddd_go/second_ddd_go/config"
	"github.com/johnnrails/ddd_go/second_ddd_go/domain"
	"github.com/johnnrails/ddd_go/second_ddd_go/infra/repositories"
	"gorm.io/gorm"
)

type NewsRepositoryImpl struct {
	DB *gorm.DB
}

func CreateNewsRepository() (repositories.NewsRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	return &NewsRepositoryImpl{DB: db}, nil
}

func (r *NewsRepositoryImpl) GetByID(id int) (*domain.News, error) {
	news := &domain.News{}
	if err := r.DB.Preload("Topic").First(&news, id).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *NewsRepositoryImpl) GetAll() ([]domain.News, error) {
	news := []domain.News{}
	if err := r.DB.Preload("Topic").Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *NewsRepositoryImpl) Save(news *domain.News) error {
	if err := r.DB.Save(&news).Error; err != nil {
		return err
	}
	return nil
}

func (r *NewsRepositoryImpl) Remove(id int) error {
	transaction := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	if err := transaction.Delete(id).Error; err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit().Error
}

func (r *NewsRepositoryImpl) Update(id int, new domain.News) error {
	if err := r.DB.Model(id).UpdateColumns(new).Error; err != nil {
		return err
	}
	return nil
}

func (r *NewsRepositoryImpl) GetByStatus(status string) ([]domain.News, error) {
	news := []domain.News{}
	if err := r.DB.Where("status = ?", status).Preload("Topic").Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *NewsRepositoryImpl) GetBySlug(slug string) ([]*domain.News, error) {
	rows, _ := r.DB.Raw("SELECT * FROM news WHERE slug = ?)", slug).Rows()
	defer rows.Close()

	news := make([]*domain.News, 0)

	for rows.Next() {
		n := &domain.News{}
		err := rows.Scan(&n.ID, &n.Title, &n.Slug, &n.Content, &n.Status)
		if err != nil {
			return nil, err
		}
		news = append(news, n)
	}
	return news, nil
}
