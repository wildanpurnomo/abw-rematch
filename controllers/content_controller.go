package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

func GetUserContents(c *gin.Context) {
	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
		return
	}

	var contents []models.Content
	if err := repositories.Repo.GetContentByUserId(&contents, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": contents})
	}
}

func CreateContent(c *gin.Context) {
	var input models.CreateContentInput
	if err := c.Bind(&input); err != nil {
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
	result := repositories.Repo.GetContentByUserIdAndTitle(&content, userId, input.Title)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// create slug
		var user models.User
		if err := repositories.Repo.FetchUserById(&user, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		input.Slug = fmt.Sprintf(
			"%s-%s",
			strings.ReplaceAll(user.Username, " ", "-"),
			strings.ReplaceAll(input.Title, " ", "-"),
		)

		// begin process upload file if exists
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["media"]
			for index, file := range files {
				if ValidateUploadFileType(file.Filename) {
					bucketName := fmt.Sprintf("media-%d-%d", time.Now().Unix(), index)
					if err := libs.UploadLib.BeginUpload(file, bucketName); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}

					publicUrl := fmt.Sprintf(os.Getenv("STORAGE_PUBLIC_URL"), bucketName)
					input.MediaUrls = append(input.MediaUrls, publicUrl)
				}
			}
		}

		content := models.Content{
			Title:       input.Title,
			Description: input.Description,
			MediaUrls:   input.MediaUrls,
			YoutubeUrl:  input.YoutubeUrl,
			Slug:        input.Slug,
		}
		if err := repositories.Repo.CreateNewContent(&user, content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": content})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be unique"})
		return
	}
}
