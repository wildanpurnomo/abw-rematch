package libs

import (
	"context"
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

type CookieAccess struct {
	GinContext *gin.Context
}

var (
	AuthContextKey         = "auth"
	CookieSetterContextKey = "cookie-setter"
	JwtSecret              = []byte(os.Getenv("JWT_SECRET"))
	PublicEndpoints        = []string{
		"api/auth/login",
		"api/auth/register",
		"api/content/browse/",
		"api/gql",
	}
)

func (c *CookieAccess) SetJwtToken(token string) {
	c.GinContext.SetCookie("jwt", token, 60*60*24, "/", "", false, true)
}

func GetCookieSetter(ctx context.Context) *CookieAccess {
	return ctx.Value(CookieSetterContextKey).(*CookieAccess)
}

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

		c.Request = c.Request.WithContext(
			context.WithValue(
				context.Background(),
				CookieSetterContextKey,
				&CookieAccess{
					GinContext: c,
				},
			),
		)
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
