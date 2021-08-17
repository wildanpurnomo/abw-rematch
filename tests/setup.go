package tests

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	gqlschema "github.com/wildanpurnomo/abw-rematch/gql/schema"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

type MockSQLQuery struct {
	Query     string
	Args      []driver.Value
	Returning []*sqlmock.Rows
}

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

func BeginGraphQLServerTesting(fullPath string, cookie http.Cookie) *httptest.ResponseRecorder {
	r := InitGQLServerTesting()
	req := httptest.NewRequest(http.MethodPost, fullPath, nil)
	req.AddCookie(&cookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func StubSQLQuery(m MockSQLQuery) (*sql.DB, *gorm.DB, error) {
	sqlMockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDb, err := gorm.Open("postgres", sqlMockDb)
	if err != nil {
		sqlMockDb.Close()
		return nil, nil, err
	}

	isSelectQuery := strings.Contains(m.Query, "SELECT")
	if !isSelectQuery {
		mock.ExpectBegin()
	}

	mock.ExpectQuery(regexp.QuoteMeta(m.Query)).
		WithArgs(m.Args...).
		WillReturnRows(m.Returning...)

	if !isSelectQuery {
		mock.ExpectCommit()
	}

	repositories.InitRepository(gormDb)

	return sqlMockDb, gormDb, nil
}
