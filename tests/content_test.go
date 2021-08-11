package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wildanpurnomo/abw-rematch/controllers"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

var (
	createContentEndpoint = "/api/content/create"
)

func TestCreateContent_NoJwt(t *testing.T) {
	// init gin for testing
	r := libs.InitGinForTesting()
	r.POST(createContentEndpoint, controllers.CreateContent)

	// begin test
	form := url.Values{}
	form.Add("title", "testing")
	form.Add("body", "testing")
	req := httptest.NewRequest(http.MethodPost, createContentEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// assert status code
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// assert response body
	jsonString := w.Body.String()
	assert.Equal(t, true, strings.Contains(jsonString, `"error":"Unauthorized client"`))
}

func TestCreateContent_NoRequestPayload(t *testing.T) {
	// init gin for testing
	r := libs.InitGinForTesting()
	r.POST(createContentEndpoint, controllers.CreateContent)

	// begin test
	token, _ := libs.GenerateToken(1)
	req := httptest.NewRequest("POST", createContentEndpoint, nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// assert status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// assert response body
	jsonString := w.Body.String()
	assert.Equal(t, true, strings.Contains(jsonString, `"error":`))
}

func TestCreateContent_TitleNotUnique(t *testing.T) {
	// init gin for testing
	r := libs.InitGinForTesting()
	r.POST(createContentEndpoint, controllers.CreateContent)

	// init sql mock
	sqlMockDb, mock, _ := sqlmock.New()
	defer sqlMockDb.Close()

	gormDb, _ := gorm.Open("postgres", sqlMockDb)
	defer gormDb.Close()

	// mock query that will be executed
	const query = `SELECT * FROM "contents"  WHERE "contents"."deleted_at" IS NULL AND ((user_id = $1 AND title = $2)) ORDER BY "contents"."id" ASC LIMIT 1`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

	// assign mock db to repository
	repositories.InitRepository(gormDb)

	// create post payload
	form := url.Values{}
	form.Add("title", "testing")
	form.Add("body", "testing")

	// begin test
	token, _ := libs.GenerateToken(1)
	req := httptest.NewRequest(http.MethodPost, createContentEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// assert status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// assert response body
	jsonString := w.Body.String()
	assert.Equal(t, true, strings.Contains(jsonString, `"error":"Title must be unique"`))
}

// TODO @wildanpurnomo
func TestCreateContent_FailedUpload(t *testing.T) {
	mockUploadSvc := new(libs.MockObject)
	mockUploadSvc.On("UploadFile", mock.Anything, mock.Anything).Return(errors.New("testing error upload"))
}

// TODO @wildanpurnomo
func TestCreateContent_SuccessfulProcess(t *testing.T) {
	mockUploadSvc := new(libs.MockObject)
	mockUploadSvc.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
}
