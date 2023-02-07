package domain

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Title   string  `json:"title"`
	Slug    string  `json:"slug"`
	Content string  `json:"content"`
	Status  string  `json:"status"`
	Topic   []Topic `gorm:"many2many:news_topics;"`
}
