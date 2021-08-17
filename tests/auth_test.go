package tests

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/h2non/gock.v1"
)

func TestRegister_ValidCase(t *testing.T) {
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
	mockSqlDb, mockGormDb, err := StubSQLQuery(
		MockSQLQuery{
			Query: `INSERT INTO "users" ("username","password","profile_picture","points","unique_code","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "users"."id"`,
			Args: []driver.Value{
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
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

	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=mutation+_{register(username:"test%20username",password:"testPasssword123"){username,profile_picture,points}}`,
		http.Cookie{},
	)

	// verify response
	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"username": "test username"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"profile_picture": "Testing"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"points": 0`))
}

func TestRegister_InvalidPassword(t *testing.T) {
	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=mutation+_{register(username:"test%20username",password:"test"){username,profile_picture,points}}`,
		http.Cookie{},
	)

	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"register": null`))
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid username or password"`))
}

func TestRegister_InvalidUsername(t *testing.T) {
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=mutation+_{register(username:"test",password:"testPasssword123"){username,profile_picture,points}}`,
		http.Cookie{},
	)

	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"register": null`))
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid username or password"`))
}

func TestLogin_InvalidUsername(t *testing.T) {
	// stub sql
	mockSqlDb, mockGormDb, err := StubSQLQuery(
		MockSQLQuery{
			Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((username = $1)) ORDER BY "users"."id" ASC LIMIT 1`,
			Args:  []driver.Value{sqlmock.AnyArg()},
			Returning: []*sqlmock.Rows{
				sqlmock.
					NewRows([]string{"username"}).
					AddRow("xqcL"),
			},
		},
	)
	if err != nil {
		t.Fatalf("error init sqlmock: %v", err)
	}
	defer mockSqlDb.Close()
	defer mockGormDb.Close()

	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=mutation+_{login(username:"test%20username",password:"testPasssword123"){username,profile_picture,points}}`,
		http.Cookie{},
	)

	// verify response
	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"login": null`))
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid username or password"`))
}

func TestLogin_InvalidPassword(t *testing.T) {
	// stub sql
	hash, _ := bcrypt.GenerateFromPassword([]byte("xqcL"), bcrypt.MinCost)
	mockSqlDb, mockGormDb, err := StubSQLQuery(
		MockSQLQuery{
			Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((username = $1)) ORDER BY "users"."id" ASC LIMIT 1`,
			Args:  []driver.Value{sqlmock.AnyArg()},
			Returning: []*sqlmock.Rows{
				sqlmock.
					NewRows([]string{"username", "password"}).
					AddRow("xqcL", hash),
			},
		},
	)
	if err != nil {
		t.Fatalf("error init sqlmock: %v", err)
	}
	defer mockSqlDb.Close()
	defer mockGormDb.Close()

	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=mutation+_{login(username:"xqcL",password:"monkaW"){username,profile_picture,points}}`,
		http.Cookie{},
	)

	// verify response
	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"login": null`))
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid username or password"`))
}

func TestLogin_ValidCase(t *testing.T) {
	// stub sql
	hash, _ := bcrypt.GenerateFromPassword([]byte("monkaW"), bcrypt.MinCost)
	mockSqlDb, mockGormDb, err := StubSQLQuery(
		MockSQLQuery{
			Query: `SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((username = $1)) ORDER BY "users"."id" ASC LIMIT 1`,
			Args:  []driver.Value{sqlmock.AnyArg()},
			Returning: []*sqlmock.Rows{
				sqlmock.
					NewRows([]string{"username", "password", "profile_picture", "points"}).
					AddRow("xqcL", hash, "ayayaclap", 0),
			},
		},
	)
	if err != nil {
		t.Fatalf("error init sqlmock: %v", err)
	}
	defer mockSqlDb.Close()
	defer mockGormDb.Close()

	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=mutation+_{login(username:"xqcL",password:"monkaW"){username,profile_picture,points}}`,
		http.Cookie{},
	)

	// verify response
	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"username": "xqcL"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"profile_picture": "ayayaclap"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"points": 0`))
}

func TestAuthenticate_InvalidCase(t *testing.T) {
	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=query+_{authenticate{username,profile_picture,points}}`,
		http.Cookie{},
	)

	// verify response
	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"authenticate": null`))
	assert.Equal(t, true, strings.Contains(jsonString, `"errors":`))
	assert.Equal(t, true, strings.Contains(jsonString, `"message": "Invalid token or user not found"`))
}

func TestAuthenticate_ValidCase(t *testing.T) {
	const userId = 1

	// stub sql
	mockSqlDb, mockGormDb, err := StubSQLQuery(
		MockSQLQuery{
			Query: `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = $1)) ORDER BY "users"."id" ASC LIMIT 1`,
			Args:  []driver.Value{fmt.Sprint(userId)},
			Returning: []*sqlmock.Rows{
				sqlmock.
					NewRows([]string{"username", "profile_picture", "points"}).
					AddRow("xqcL", "pogU", 0),
			},
		},
	)
	if err != nil {
		t.Fatalf("error init sqlmock: %v", err)
	}
	defer mockSqlDb.Close()
	defer mockGormDb.Close()

	// Inject token
	token, _ := libs.GenerateToken(userId)

	// begin test
	testResponse := BeginGraphQLServerTesting(
		`/api/gql?query=query+_{authenticate{username,profile_picture,points}}`,
		http.Cookie{
			Name:  "jwt",
			Value: token,
		},
	)

	// verify response
	jsonString := testResponse.Body.String()
	assert.Equal(t, http.StatusOK, testResponse.Code)
	assert.Equal(t, true, strings.Contains(jsonString, `"username": "xqcL"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"profile_picture": "pogU"`))
	assert.Equal(t, true, strings.Contains(jsonString, `"points": 0`))
}
