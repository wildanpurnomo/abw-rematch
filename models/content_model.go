package models

import (
	"strings"
	"time"

	"github.com/lib/pq"
)

type Content struct {
	ID           uint           `json:"content_id" gorm:"primary_key"`
	UserID       uint           `json:"-" gorm:"not null"`
	Title        string         `json:"title" gorm:"not null; default:null"`
	Body         string         `json:"body" gorm:"default:null"`
	MediaUrls    pq.StringArray `json:"media_urls" gorm:"type:text[]; default:null"`
	YoutubeUrl   string         `json:"youtube_url" gorm:"default:null"`
	Slug         string         `json:"slug" gorm:"unique; not null; default:null"`
	Redirections []Redirection  `json:",omitempty" gorm:"foreignKey: New; references: Slug; constraint: OnUpdate:CASCADE, OnDelete:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `json:"-" sql:"index"`
}

func (c *Content) GetMediaBucketNames() []string {
	bktNames := []string{}
	for _, fileUrl := range c.MediaUrls {
		lastIndex := strings.Index(fileUrl, "?")
		bktName := fileUrl[70:lastIndex]
		bktNames = append(bktNames, bktName)
	}
	return bktNames
}

type CreateContentInput struct {
	Title      string         `json:"title" form:"title" binding:"required"`
	Body       string         `json:"body" form:"body" binding:"required"`
	MediaUrls  pq.StringArray `json:"media_urls"`
	YoutubeUrl string         `json:"youtube_url"`
	Slug       string         `json:"slug"`
}
