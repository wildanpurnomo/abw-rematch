package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func VerifyPassword(hashed []byte, plain []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashed, plain); err != nil {
		return false
	}

	return true
}

func VerifyJwt(c *gin.Context) (uint, bool) {
	cookie, err := c.Request.Cookie("jwt")
	if err != nil {
		return 0, false
	}

	cookieValue := cookie.Value
	claims := &models.JwtClaims{}

	token, err := jwt.ParseWithClaims(cookieValue, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			fmt.Print(err)
		}
		return 0, false
	}

	return claims.UserID, true
}

func GenerateToken(userId uint) (string, bool) {
	claims := &models.JwtClaims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := sign.SignedString(jwtSecret)
	if err != nil {
		return "", false
	}

	return token, true
}

func Authenticate(c *gin.Context) {
	userId, status := VerifyJwt(c)
	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": false})
		return
	}

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
	if !VerifyPassword([]byte(user.Password), []byte(input.Password)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// invoke token
	token, status := GenerateToken(user.ID)
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
	if !ValidateUsername(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 8 characters long"})
		return
	}

	// password validation
	if !ValidatePassword(input.Password) {
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
	token, status := GenerateToken(newUser.ID)
	if !status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	c.SetCookie("jwt", token, 60*60*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}
