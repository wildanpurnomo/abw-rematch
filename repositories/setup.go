package repositories

import "github.com/jinzhu/gorm"

type Repository struct {
	db *gorm.DB
}

var Repo *Repository

func InitRepository(db *gorm.DB) {
	Repo = &Repository{db: db}
}
