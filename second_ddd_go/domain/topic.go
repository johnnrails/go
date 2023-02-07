package domain

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	Name string `json:"name"`
	Slug string `json:"slug"`
	News []News `gorm:"many2many:news_topics;"`
}
