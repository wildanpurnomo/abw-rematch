package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
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

	// trim title and body
	input.Title = strings.TrimSpace(input.Title)
	input.Body = strings.TrimSpace(input.Body)

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
		input.Slug = slug.Make(fmt.Sprintf("%s %s", user.UniqueCode, input.Title))

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
			Title:      input.Title,
			Body:       input.Body,
			MediaUrls:  input.MediaUrls,
			YoutubeUrl: input.YoutubeUrl,
			Slug:       input.Slug,
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

func UpdateContent(c *gin.Context) {
	// verify that input is valid form-data
	var input models.CreateContentInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// verify jwt
	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
		return
	}

	// extract contentId from path param
	contentId, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contentId"})
		return
	}

	// trim title and body
	input.Title = strings.TrimSpace(input.Title)
	input.Body = strings.TrimSpace(input.Body)

	// check whether user is authorized to edit requested content
	var content models.Content
	if err := repositories.Repo.GetContentByUserIdAndContentId(&content, userId, uint(contentId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content not found or unauthorized user"})
		return
	}

	// check whether media payload is exist
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["media"]
		if len(files) > 0 {
			// delete existing media
			for _, fileUrl := range content.MediaUrls {
				lastIndex := strings.Index(fileUrl, "?")
				bktName := fileUrl[69:lastIndex]
				libs.UploadLib.BeginDeleteFile(bktName)
			}

			// upload new media and assign new urls
			for index, file := range files {
				if ValidateUploadFileType(file.Filename) {
					bucketName := fmt.Sprintf("media-%d-%d", time.Now().Unix(), index)
					if err := libs.UploadLib.BeginUpload(file, bucketName); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "File processing failed"})
						return
					}

					publicUrl := fmt.Sprintf(os.Getenv("STORAGE_PUBLIC_URL"), bucketName)
					input.MediaUrls = append(input.MediaUrls, publicUrl)
				}
			}
		}
	}

	// check whether title is changed
	oldSlug := content.Slug
	if input.Title != content.Title {
		// check whether title is unique
		contentResult := repositories.Repo.GetContentByUserIdAndTitle(&models.Content{}, userId, input.Title)
		if errors.Is(contentResult.Error, gorm.ErrRecordNotFound) {
			input.Slug = slug.Make(fmt.Sprintf("%s-%s", content.Slug[:10], input.Title))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"data": "Invalid title"})
			return
		}
	}

	// execute update query
	if err := repositories.Repo.UpdateContent(
		&content,
		models.Content{
			Title:      input.Title,
			Body:       input.Body,
			MediaUrls:  input.MediaUrls,
			YoutubeUrl: input.YoutubeUrl,
			Slug:       input.Slug,
		},
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update requested content"})
	} else {
		// create new redirection record
		if oldSlug != content.Slug {
			repositories.Repo.CreateNewRedirection(
				&models.Redirection{
					Old: oldSlug,
					New: input.Slug,
				},
			)
		}

		c.JSON(http.StatusOK, gin.H{"data": content})
	}
}
