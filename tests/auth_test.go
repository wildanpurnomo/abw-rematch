package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
	"github.com/wildanpurnomo/abw-rematch/controllers"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	"gopkg.in/h2non/gock.v1"
)

func TestLogin_NoJSONPayload(t *testing.T) {
	// init gin
	r := libs.InitGinForTesting()
	r.POST("/api/auth/login", controllers.Login)

	// begin test
	req := httptest.NewRequest("POST", "/api/auth/login", nil) // no json payload
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// assert status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// assert response body
	assert.Equal(t, true, strings.Contains(w.Body.String(), `"error":"Invalid username or password"`))
}

func TestRegister_NoJSONPayload(t *testing.T) {
	// init gin
	r := libs.InitGinForTesting()
	r.POST("/api/auth/register", controllers.Register)

	// begin test
	req := httptest.NewRequest("POST", "/api/auth/register", nil) // no json payload
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// assert status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// assert response body
	assert.Equal(t, true, strings.Contains(w.Body.String(), `"error":"EOF"`))
}

func TestRegister_ValidCase(t *testing.T) {
	// Init gin
	r := libs.InitGinForTesting()
	r.POST("/api/auth/register", controllers.Register)

	// init sqlmock
	sqlMockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error init sqlmock: %v", err)
	}
	defer sqlMockDb.Close()

	gormDb, err := gorm.Open("postgres", sqlMockDb)
	if err != nil {
		t.Fatalf("error init gormDb: %v", err)
	}
	defer gormDb.Close()

	// init gock to mock randomUserApi call
	defer gock.Off()

	mockResults := []models.RandomUser{
		{
			ProfilePicture: models.ProfilePicture{Medium: "Testing"},
		},
	}
	mockRandomUserAPIResponse := models.RandomUserAPIResponse{
		Results: mockResults,
	}
	gock.New("https://randomuser.me/api").Get("/").Reply(200).JSON(mockRandomUserAPIResponse)

	// mock insert SQL
	const sqlInsert = `INSERT INTO "users" ("username","password","profile_picture","points","unique_code","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "users"."id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// assign mock sql db to repository
	repositories.InitRepository(gormDb)

	// create POST payload
	authInput := models.UserAuthInput{
		Username: "test username",
		Password: "testPassword123",
	}
	jsonTest, _ := json.Marshal(authInput)

	// begin test
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonTest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// verify http status
	assert.Equal(t, http.StatusOK, w.Code)

	// verify response body
	jsonString := w.Body.String()
	assert.Equal(t, true, strings.Contains(jsonString, `"username":"test username"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"profile_picture":"Testing"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"points":0`))
	assert.Equal(t, false, strings.Contains(jsonString, `"Password":`))
	assert.Equal(t, false, strings.Contains(jsonString, `"UniqueCode":`))
	assert.Equal(t, false, strings.Contains(jsonString, `"ID":`))
}

func TestRegister_InvalidPassword(t *testing.T) {
	r := libs.InitGinForTesting()
	r.POST("/api/auth/register", controllers.Register)

	authInput := models.UserAuthInput{
		Username: "test username",
		Password: "test",
	}
	jsonTest, _ := json.Marshal(authInput)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonTest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	jsonString := w.Body.String()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"error":"Password must be at least 8 characters long, contains min 1 uppercase, min 1 lowercase and 1 number"`))
}

func TestRegister_InvalidUsername(t *testing.T) {
	r := libs.InitGinForTesting()
	r.POST("/api/auth/register", controllers.Register)

	authInput := models.UserAuthInput{
		Username: "test",
		Password: "testPassword123",
	}

	jsonTest, _ := json.Marshal(authInput)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonTest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	jsonString := w.Body.String()
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"error":"Username must be at least 8 characters long"`))
}
