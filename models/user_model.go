package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	ID             uint      `json:"-" gorm:"primary_key"`
	Username       string    `json:"username" gorm:"unique"`
	Password       string    `json:"-"`
	ProfilePicture string    `json:"profile_picture"`
	Points         int       `json:"points"`
	UniqueCode     string    `json:"-" gorm:"unique; not null; default:null"`
	Contents       []Content `json:",omitempty"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `json:"-" sql:"index"`
}

func (u *User) Name() string {
	return "UserModelType"
}

func (u *User) Description() string {
	return "User's data"
}

func (u *User) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(out)
}

func (u *User) Error() error {
	return errors.New("Whoops user")
}

type UserAuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUsernameInput struct {
	Username string `json:"username" binding:"required"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type RandomUserAPIResponse struct {
	Results []RandomUser `json:"results"`
}

type RandomUser struct {
	ProfilePicture ProfilePicture `json:"picture"`
}

type ProfilePicture struct {
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Thumbnail string `json:"thumbnail"`
}

type JwtClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}
