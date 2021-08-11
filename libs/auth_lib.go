package libs

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wildanpurnomo/abw-rematch/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	AuthContextKey  = "auth"
	JwtSecret       = []byte(os.Getenv("JWT_SECRET"))
	PublicEndpoints = []string{
		"api/auth/login",
		"api/auth/register",
		"api/content/browse/",
	}
)

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
		return JwtSecret, nil
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
	token, err := sign.SignedString(JwtSecret)
	if err != nil {
		return "", false
	}

	return token, true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsPublicEndpoint(c.Request.URL.Path) {
			userId, status := VerifyJwt(c)
			if !status {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized client"})
				c.Abort()
				return
			}

			c.Set(AuthContextKey, fmt.Sprint(userId))
		}
		c.Next()
	}
}

func IsPublicEndpoint(path string) bool {
	for _, endpoint := range PublicEndpoints {
		if strings.Contains(path, endpoint) {
			return true
		}
	}

	return false
}
