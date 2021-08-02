package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Content struct {
	UserID      uint           `json:"user_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null;default:null"`
	Description string         `json:"description"`
	MediaUrls   pq.StringArray `json:"media_urls" gorm:"type:text[]"`
	YoutubeUrl  string         `json:"youtube_url"`
	Slug        string         `json:"slug" gorm:"unique;not null;default:null"`
	gorm.Model
}

type CreateContentInput struct {
	Title       string         `json:"title" form:"title" binding:"required"`
	Description string         `json:"description" form:"description" binding:"required"`
	MediaUrls   pq.StringArray `json:"media_urls"`
	YoutubeUrl  string         `json:"youtube_url"`
	Slug        string         `json:"slug"`
}
