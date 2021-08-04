package models

import "time"

// used to record changes of slug/url
// e.g:
// content slug will change everytime content title is updated.
// when client requests a content by old slug, this table will help server to redirect client to updated slug
type Redirection struct {
	ID        uint       `json:"-" gorm:"primary_key"`
	Old       string     `json:"-"`
	New       string     `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
