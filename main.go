package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	gqlschema "github.com/wildanpurnomo/abw-rematch/gql/schema"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	controllers "github.com/wildanpurnomo/abw-rematch/rest-controllers"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	if os.Getenv("ENV_SCHEMA") == "https" {
		gin.SetMode(gin.ReleaseMode)
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
	r.Use(libs.AuthMiddleware())

	graphqlSchema, err := gqlschema.InitSchema()
	if err != nil {
		log.Fatal(err)
	}
	gqlHandler := handler.New(&handler.Config{
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: true,
	})
	gqlHandlerFunc := gin.HandlerFunc(func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	gqlRoutes := r.Group("/api")
	{
		gqlRoutes.GET("/gql", gqlHandlerFunc)
		gqlRoutes.POST("/gql", gqlHandlerFunc)
	}

	contentRestRoutes := r.Group("/api")
	{
		contentRestRoutes.GET("/content/browse/:slug", controllers.GetContentBySlug)
		contentRestRoutes.POST("/content/create", controllers.CreateContent)
		contentRestRoutes.PUT("/content/update/:contentId", controllers.UpdateContent)
	}

	r.Run()
}
