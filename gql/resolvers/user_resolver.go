package gqlresolvers

import (
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/wildanpurnomo/abw-rematch/libs"
	"github.com/wildanpurnomo/abw-rematch/models"
	"github.com/wildanpurnomo/abw-rematch/repositories"
)

var UpdateUsernameResolver = func(params graphql.ResolveParams) (interface{}, error) {
	cookieAccess := libs.GetContextValues(params.Context)
	userId := cookieAccess.UserID
	if userId == "0" {
		return nil, errors.New("Invalid token or user not found")
	}

	var input models.UpdateUsernameInput

	// trim username
	input.Username = strings.TrimSpace(params.Args["username"].(string))

	// username validation
	if !libs.ValidateUsername(input.Username) {
		return nil, errors.New("Invalid new username")
	}

	// fetch user from DB
	var user models.User
	if err := repositories.Repo.FetchUserById(&user, userId); err != nil {
		return nil, errors.New("Invalid token or user not found")
	}

	// update username
	update := user
	update.Username = input.Username
	if err := repositories.Repo.UpdateUser(&user, update); err != nil {
		return nil, errors.New("Whoops!")
	}

	return user, nil
}
