package gqlresolvers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/graphql-go/graphql"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
	"golang.org/x/crypto/bcrypt"
)

var LoginResolver = func(params graphql.ResolveParams) (interface{}, error) {
	var input models.UserAuthInput

	// populate as input object
	input.Username = params.Args["username"].(string)
	input.Password = params.Args["password"].(string)

	var user models.User
	if err := repositories.Repo.FetchUserByUsername(&user, input.Username); err != nil {
		return nil, errors.New("Invalid username or password")
	}

	// verify password
	if !libs.VerifyPassword([]byte(user.Password), []byte(input.Password)) {
		return nil, errors.New("Invalid username or password")
	}

	// invoke token
	token, status := libs.GenerateToken(user.ID)
	if !status {
		return nil, errors.New("Invalid username or password")
	}

	cookieAccess := libs.GetCookieSetter(params.Context)
	cookieAccess.SetJwtToken(token)

	return user, nil
}

var RegisterResolver = func(params graphql.ResolveParams) (interface{}, error) {
	var input models.UserAuthInput

	// populate as input object
	input.Username = strings.TrimSpace(params.Args["username"].(string))
	input.Password = params.Args["password"].(string)

	// username validation
	if !libs.ValidateUsername(input.Username) {
		return nil, errors.New("Invalid username or password")
	}

	// password validation
	if !libs.ValidatePassword(input.Password) {
		return nil, errors.New("Invalid username or password")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return nil, errors.New("Whoops!")
	}

	// get profile picture
	res, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		return nil, errors.New("Whoops!")
	}
	defer res.Body.Close()
	var randomUserApiResponse models.RandomUserAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&randomUserApiResponse); err != nil {
		return nil, errors.New("Whoops!")
	}

	// save to db
	newUser := models.User{
		Username:       input.Username,
		Password:       string(hash),
		ProfilePicture: randomUserApiResponse.Results[0].ProfilePicture.Medium,
		Points:         0,
		UniqueCode:     uniuri.NewLen(10),
	}
	if err := repositories.Repo.CreateNewUser(&newUser); err != nil {
		return nil, errors.New("Whoops!")
	}

	// invoke token
	token, status := libs.GenerateToken(newUser.ID)
	if !status {
		return nil, errors.New("Whoops!")
	}
	cookieAccess := libs.GetCookieSetter(params.Context)
	cookieAccess.SetJwtToken(token)

	return newUser, nil
}
