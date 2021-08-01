package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Username       string `json:"username" gorm:"unique"`
	Password       string `json:"password,omitempty"`
	ProfilePicture string `json:"profile_picture"`
	Points         int    `json:"points"`
	Contents       []Content
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
