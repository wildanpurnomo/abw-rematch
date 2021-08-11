package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(c *gin.Context) {
	userId := c.GetString(libs.AuthContextKey)

	var user models.User
	if err := repositories.Repo.FetchUserById(&user, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func Login(c *gin.Context) {
	var input models.UserAuthInput

	// json validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	var user models.User
	if err := repositories.Repo.FetchUserByUsername(&user, input.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// verify password
	if !libs.VerifyPassword([]byte(user.Password), []byte(input.Password)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// invoke token
	token, status := libs.GenerateToken(user.ID)
	if !status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	c.SetCookie("jwt", token, 60*60*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Register(c *gin.Context) {
	var input models.UserAuthInput

	// json validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trim username
	input.Username = strings.TrimSpace(input.Username)

	// username validation
	if !libs.ValidateUsername(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 8 characters long"})
		return
	}

	// password validation
	if !libs.ValidatePassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long, contains min 1 uppercase, min 1 lowercase and 1 number"})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get profile picture
	res, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer res.Body.Close()
	var randomUserApiResponse models.RandomUserAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&randomUserApiResponse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// save to db
	newUser := models.User{
		Username:       input.Username,
		Password:       string(hash),
		ProfilePicture: randomUserApiResponse.Results[0].ProfilePicture.Medium,
		Points:         0,
		UniqueCode:     uniuri.NewLen(10),
	}
	if err := repositories.Repo.CreateNewUser(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// invoke token
	token, status := libs.GenerateToken(newUser.ID)
	if !status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	c.SetCookie("jwt", token, 60*60*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
