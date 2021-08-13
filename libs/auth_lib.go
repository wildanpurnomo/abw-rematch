package libs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wildanpurnomo/abw-rematch/models"
	"golang.org/x/crypto/bcrypt"
)

type ContextValues struct {
	GinContext *gin.Context
	UserID     string
}

var (
	AuthContextKey  = "auth"
	ContextValueKey = "context-values"
	JwtSecret       = []byte(os.Getenv("JWT_SECRET"))
)

func (c *ContextValues) InvalidateToken() {
	c.GinContext.SetCookie("jwt", "", 1, "/", "", false, true)
}

func (c *ContextValues) SetJwtToken(token string) {
	c.GinContext.SetCookie("jwt", token, 60*60*24, "/", "", false, true)
}

func GetContextValues(ctx context.Context) *ContextValues {
	return ctx.Value(ContextValueKey).(*ContextValues)
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
		userId, _ := VerifyJwt(c)
		c.Request = c.Request.WithContext(
			context.WithValue(
				c.Request.Context(),
				ContextValueKey,
				&ContextValues{
					GinContext: c,
					UserID:     fmt.Sprint(userId),
				},
			),
		)
		c.Next()
	}
}
