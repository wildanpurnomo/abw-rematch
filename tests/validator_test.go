package tests

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/wildanpurnomo/abw-rematch/controllers"
)

func TestValidateUsername_ValidCase(t *testing.T) {
	testCase := "this is valid username"

	assert.Equal(t, true, controllers.ValidateUsername(testCase))
}

func TestValidateUsername_InvalidCase(t *testing.T) {
	testCase := "tooShort"

	assert.Equal(t, false, controllers.ValidateUsername(testCase))
}

func TestValidatePassword_ValidCase(t *testing.T) {
	testCases := []string{
		"thisIsCorrect123",
		"thisOneAlso123",
		"testPassword123",
	}

	for _, testCase := range testCases {
		assert.Equal(t, true, controllers.ValidatePassword(testCase))
	}
}

func TestValidatePassword_InvalidCase(t *testing.T) {
	testCases := []string{
		"thisWrong",
		"thisIsWrong",
		"Wrong1",
	}

	for _, testCase := range testCases {
		assert.Equal(t, false, controllers.ValidatePassword(testCase))
	}
}
