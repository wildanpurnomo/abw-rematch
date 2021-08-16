package tests

import (
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/stretchr/testify/mock"
	gqlschema "github.com/wildanpurnomo/abw-rematch/gql/schema"
	"github.com/wildanpurnomo/abw-rematch/libs"
)

type MockObject struct {
	mock.Mock
}

func (s *MockObject) UploadFile(fileHeader *multipart.FileHeader, bucketName string) error {
	return nil
}

func InitRESTServerTesting() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(libs.AuthMiddleware())
	return router
}

func InitGQLServerTesting() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
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

	return r
}
