package tests

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	gqlschema "github.com/wildanpurnomo/abw-rematch/gql/schema"
	"github.com/wildanpurnomo/abw-rematch/libs"
)

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
