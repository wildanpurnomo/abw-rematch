package tests

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	"gopkg.in/h2non/gock.v1"
)

func TestRegister_ValidCase(t *testing.T) {
	// Init gin
	r := InitGQLServerTesting()

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

	// begin test
	req := httptest.NewRequest(
		http.MethodPost,
		`/api/gql?query=mutation+_{register(username:"test%20username",password:"testPasssword123"){username,profile_picture,points}}`,
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// verify http status
	assert.Equal(t, http.StatusOK, w.Code)

	// verify response body
	jsonString := w.Body.String()
	assert.Equal(t, true, strings.Contains(jsonString, `"username": "test username"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"profile_picture": "Testing"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"points": 0`))
}

func TestRegister_InvalidPassword(t *testing.T) {
	r := InitGQLServerTesting()

	req := httptest.NewRequest(
		http.MethodPost,
		`/api/gql?query=mutation+_{register(username:"test%20username",password:"test"){username,profile_picture,points}}`,
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	jsonString := w.Body.String()
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid username or password"`))
}

func TestRegister_InvalidUsername(t *testing.T) {
	r := InitGQLServerTesting()

	req := httptest.NewRequest(
		http.MethodPost,
		`/api/gql?query=mutation+_{register(username:"test",password:"testPasssword123"){username,profile_picture,points}}`,
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	jsonString := w.Body.String()
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid username or password"`))
}
