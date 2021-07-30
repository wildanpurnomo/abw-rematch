package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wildanpurnomo/abw-rematch/models"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func ConnectFirebase() {
	opt := option.WithCredentialsFile("firebaseServiceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error init firebase: %v", err)
	}

	firebaseApp = app
}

func GetUserContents(c *gin.Context) {
	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
		return
	}

	var contents []models.Content
	models.DB.Where("user_id = ?", userId).Find(&contents)

	c.JSON(http.StatusOK, gin.H{"data": contents})
}

func CreateContent(c *gin.Context) {
	var input models.CreateContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
		return
	}

	// trim title and description
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	// verify title is unique
	var content models.Content
	result := models.DB.Where("user_id = ? AND title = ?", userId, input.Title).First(&content)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var user models.User
		if err := models.DB.Where("id = ?", userId).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// check if request contains file
		form, _ := c.MultipartForm()
		files := form.File["media"]
		if len(files) > 0 {
			// upload files to gcs
			_, err := firebaseApp.Storage(context.Background())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "There's a problem with upload service"})
				return
			}
		}

		appended := models.DB.Model(&user).Association("Contents").Append(&models.Content{Title: input.Title, Description: input.Description, MediaUrls: input.MediaUrls, YoutubeUrl: input.YoutubeUrl})
		c.JSON(http.StatusOK, gin.H{"data": appended})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be unique"})
		return
	}
}
