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

var LogoutResolver = func(params graphql.ResolveParams) (interface{}, error) {
	contextValue := libs.GetContextValues(params.Context)
	contextValue.InvalidateToken()
	return true, nil
}

var AuthenticateResolver = func(params graphql.ResolveParams) (interface{}, error) {
	contextValue := libs.GetContextValues(params.Context)
	userId := contextValue.UserID
	if userId == "0" {
		return nil, errors.New("Invalid token or user not found")
	}

	var user models.User
	if err := repositories.Repo.FetchUserById(&user, userId); err != nil {
		return nil, errors.New("Invalid token or user not found")
	}

	return user, nil
}

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

	contextValue := libs.GetContextValues(params.Context)
	contextValue.SetJwtToken(token)

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
	contextValue := libs.GetContextValues(params.Context)
	contextValue.SetJwtToken(token)

	return newUser, nil
}
