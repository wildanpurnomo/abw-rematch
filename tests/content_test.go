package tests

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wildanpurnomo/abw-rematch/libs"
	controllers "github.com/wildanpurnomo/abw-rematch/rest-controllers"
)

var (
	createContentEndpoint = "/api/content/create"
)

func TestCreateContent_NoJwt(t *testing.T) {
	// init gin for testing
	r := InitRESTServerTesting()
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
	assert.Equal(t, true, strings.Contains(jsonString, `"error":"Invalid token or user not found"`))
}

func TestCreateContent_NoRequestPayload(t *testing.T) {
	// init gin for testing
	r := InitRESTServerTesting()
	r.POST(createContentEndpoint, controllers.CreateContent)

	// begin test
	token, _ := libs.GenerateToken(1)
	req := httptest.NewRequest(http.MethodPost, createContentEndpoint, nil)
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
	r := InitRESTServerTesting()
	r.POST(createContentEndpoint, controllers.CreateContent)

	// stub sql
	mockSqlDb, mockGormDb, err := StubSQLQuery(
		MockSQLQuery{
			Query: `SELECT * FROM "contents"  WHERE "contents"."deleted_at" IS NULL AND ((user_id = $1 AND title = $2)) ORDER BY "contents"."id" ASC LIMIT 1`,
			Args: []driver.Value{
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			},
			Returning: []*sqlmock.Rows{
				sqlmock.NewRows([]string{"id"}).AddRow(1),
			},
		},
	)
	if err != nil {
		t.Fatalf("error mock sql: %v", err)
	}
	defer mockSqlDb.Close()
	defer mockGormDb.Close()

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
	mockUploadSvc := new(MockObject)
	mockUploadSvc.On("UploadFile", mock.Anything, mock.Anything).Return(errors.New("testing error upload"))
}

// TODO @wildanpurnomo
func TestCreateContent_SuccessfulProcess(t *testing.T) {
	mockUploadSvc := new(MockObject)
	mockUploadSvc.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
}
