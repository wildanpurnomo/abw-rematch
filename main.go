package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wildanpurnomo/abw-rematch/controllers"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
	r := gin.Default()

	dbConn, err := models.ConnectDatabase()
	if err != nil {
		fmt.Printf("error init DB: %v", err)
		return
	}

	repositories.InitRepository(dbConn)

	firebaseApp, err := libs.ConnectFirebase()
	if err != nil {
		fmt.Printf("error init firebaseApp: %v", err)
		return
	}

	libs.InitUploadLib(libs.UploadService{App: firebaseApp})

	r.Use(libs.CORSMiddleware())

	authRoutes := r.Group("auth")
	{
		authRoutes.GET("/authenticate", controllers.Authenticate)
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/logout", controllers.Logout)
	}

	userRoutes := r.Group("user")
	{
		userRoutes.PUT("/update-username", controllers.UpdateUsername)
		userRoutes.PUT("/update-password", controllers.UpdatePassword)
	}

	contentRoutes := r.Group("content")
	{
		contentRoutes.GET("/me", controllers.GetUserContents)
		contentRoutes.POST("/create", controllers.CreateContent)
	}

	r.Run()
}
