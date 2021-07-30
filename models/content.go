package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Content struct {
	UserID      uint           `json:"user_id" gorm:"not null"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	MediaUrls   pq.StringArray `json:"media_urls" gorm:"type:text[]"`
	YoutubeUrl  string         `json:"youtube_url"`
	gorm.Model
}

type CreateContentInput struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	MediaUrls   pq.StringArray `json:"media_urls"`
	YoutubeUrl  string         `json:"youtube_url"`
}
