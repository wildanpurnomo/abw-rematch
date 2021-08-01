package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wildanpurnomo/abw-rematch/models"
	"golang.org/x/crypto/bcrypt"
)

func UpdatePassword(c *gin.Context) {
	var input models.UpdatePasswordInput

	// json validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password or new password"})
		return
	}

	// jwt validation
	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
		return
	}

	// fetch user from DB
	var user models.User
	if err := models.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// verify password
	if !VerifyPassword([]byte(user.Password), []byte(input.OldPassword)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password or new password"})
		return
	}

	// password validation
	if !ValidatePassword(input.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long, contains min 1 uppercase, min 1 lowercase and 1 number"})
		return
	}

	// hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update password
	update := user
	update.Password = string(hash)
	models.DB.Model(&user).Updates(update)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateUsername(c *gin.Context) {
	var input models.UpdateUsernameInput

	// json validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		return
	}

	// jwt validation
	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
		return
	}

	// trim username
	input.Username = strings.TrimSpace(input.Username)

	// username validation
	if !ValidateUsername(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 8 characters long"})
		return
	}

	// fetch user from DB
	var user models.User
	if err := models.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// update username
	models.DB.Model(&user).Updates(input)

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": user})
}
